package docker

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/builder"
)

const (
	buildCtxDir  = "/tmp/unweave/buildctx"
	buildLogsDir = "/tmp/unweave/logs"
)

var (
	// ErrBuildFailed is returned when a build fails.
	ErrBuildFailed = &types.Error{
		Code:       http.StatusBadRequest,
		Message:    "Build failed - check the logs for more information",
		Suggestion: "Make sure your Dockerfile is valid",
	}
)

// buildImage builds an image with the Dockerfile in the given buildPath directory. It
// expects the context directory to have a Dockerfile. It will return an error if none is
// found.
//
// We might want to convert this to user the Docker SDK.
func buildImage(ctx context.Context, buildPath, image, cache string) (
	logsch chan string, errch chan error, err error,
) {
	if _, err := os.Stat(buildPath); os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("buildPath %q does not exist: %w", buildPath, err)
	}

	c := []string{
		"docker",
		"build",
		//"--cache-from", cache, // TODO
		"--build-arg", "BUILDKIT_INLINE_CACHE=1",
		"-t", image,
		buildPath,
	}
	log.Ctx(ctx).Info().Msgf("Executing command: %v", c)

	cmd := exec.CommandContext(ctx, c[0], c[1:]...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "DOCKER_BUILDKIT=1")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	errch = make(chan error)
	logsch = make(chan string, 1000) // buffer channel in case i/o is slow

	readStdout := func(stdout io.ReadCloser, output chan string) {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			output <- scanner.Text()
		}
	}
	readStderr := func(stderr io.ReadCloser, output chan string) {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			output <- scanner.Text()
		}
	}

	go readStdout(stdout, logsch)
	go readStderr(stderr, logsch)

	go func() {
		defer close(logsch)
		if e := cmd.Start(); e != nil {
			errch <- e
			return
		}
		if e := cmd.Wait(); e != nil {
			if e, ok := e.(*exec.ExitError); ok {
				// This is the build failing not the command. i.e. not Unweave's fault, so
				// we write it to the build logs and return a 400 to indicate user error.
				logsch <- e.Error()
				logsch <- fmt.Sprintf("Exit code: %d", e.ExitCode())
				errch <- ErrBuildFailed
				return
			}
			errch <- e
		}
	}()

	return logsch, errch, nil
}

func findImage(ctx context.Context, image string) (string, error) {
	cmd := exec.CommandContext(
		ctx,
		"docker",
		"image",
		"ls",
		"-q",
		image,
	)
	data, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	imageID := ""
	// Parse the output and return the first image ID. Use line break as the delimiter.
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		imageID = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error parsing image ls output: %w", err)
	}
	if imageID == "" {
		return "", fmt.Errorf("image with name %q not found", image)
	}
	return imageID, nil
}

// pushImage pushes the image to the registry
func pushImage(ctx context.Context, uri string) (output string, err error) {
	cmd := exec.CommandContext(
		ctx,
		"docker",
		"push",
		uri,
	)
	data, err := cmd.CombinedOutput()
	return string(data), err
}

// saveContext will use a zip reader to parse the context bytes and save the files
// to disk in the given saveDir path.
func saveContext(saveDir string, context []byte) error {
	foundDockerfile := false
	reader := bytes.NewReader(context)

	zr, err := zip.NewReader(reader, int64(len(context)))
	if err != nil {
		return err
	}

	// Remove the directory if it already exists - should never happen
	if err := os.RemoveAll(saveDir); err != nil {
		return err
	}
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return err
	}

	for _, zipFile := range zr.File {
		if zipFile.Name == "Dockerfile" {
			foundDockerfile = true
		}

		if zipFile.FileInfo().IsDir() {
			if err := os.MkdirAll(filepath.Join(saveDir, zipFile.Name), 0755); err != nil {
				return err
			}
			continue
		}

		f, err := os.Create(filepath.Join(saveDir, zipFile.Name))
		if err != nil {
			return err
		}
		r, err := zipFile.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(f, r)
		if err != nil {
			return err
		}
	}

	if !foundDockerfile {
		return &types.Error{
			Code:       http.StatusBadRequest,
			Message:    "Dockerfile not found in context",
			Suggestion: "Make sure your Dockerfile is in the root of your context",
			Err:        fmt.Errorf("dockerfile not found in context"),
		}
	}
	return nil
}

// tagImage will tag the target image as the source image e.g. tag a cache image as the source image
func tagImage(ctx context.Context, source, target string) (string, error) {
	cmd := exec.CommandContext(
		ctx,
		"docker",
		"tag",
		source,
		target,
	)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "DOCKER_BUILDKIT=1")

	data, err := cmd.CombinedOutput()
	return string(data), err
}

// Builder is a Docker builder that implements the builder.Builder interface.
type Builder struct {
	logger      builder.LogDriver
	registryURI string
}

func (b *Builder) GetBuilder() string {
	return "docker"
}

func (b *Builder) GetImageURI(ctx context.Context, buildID, namespace, reponame string) string {
	return fmt.Sprintf("%s/%s/%s:%s", b.registryURI, namespace, reponame, buildID)
}

func (b *Builder) Logs(ctx context.Context, buildID string) ([]types.LogEntry, error) {
	ctx = log.With().Str("builder", b.GetBuilder()).Str("buildID", buildID).Logger().WithContext(ctx)
	log.Ctx(ctx).Info().Msg("Executing logs request")
	return b.logger.GetLogs(ctx, buildID)
}

func (b *Builder) Build(ctx context.Context, buildID string, buildCtx io.Reader) error {
	ctx = log.With().
		Str("builder", b.GetBuilder()).
		Str("buildID", buildID).
		Logger().WithContext(ctx)

	log.Ctx(ctx).Info().Msg("Executing build request")

	dir := filepath.Join(buildCtxDir, buildID)
	buildBytes, err := io.ReadAll(buildCtx)
	if err != nil {
		return fmt.Errorf("failed to read build context: %w", err)
	}

	if err := saveContext(dir, buildBytes); err != nil {
		return fmt.Errorf("failed to save build context: %w", err)
	}

	imageURI := fmt.Sprintf("uw-provisional:%s", buildID) // until tagged
	logsch, errch, err := buildImage(ctx, dir, imageURI, "")
	if err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}
	log.Ctx(ctx).Info().Msg("Started image build with build context at " + dir)

	var logs []types.LogEntry

	saveCtx := context.Background()
	saveCtx = log.Logger.WithContext(ctx)
	defer func() {
		log.Ctx(saveCtx).Info().Msg("Saving logs")
		if err := b.logger.SaveLogs(saveCtx, buildID, logs); err != nil {
			log.Ctx(saveCtx).Error().Err(err).Msg("Failed to save logs")
		}
	}()

	for {
		select {
		case <-saveCtx.Done():
			return nil
		case l, ok := <-logsch:
			if !ok {
				return nil
			}
			logs = append(logs, types.LogEntry{TimeStamp: time.Now(), Message: l})
		case e := <-errch:
			return e
		}
	}
}

func (b *Builder) Push(ctx context.Context, buildID, namespace, reponame string) error {
	ctx = log.With().
		Str("builder", b.GetBuilder()).
		Str("buildID", buildID).
		Logger().WithContext(ctx)

	log.Ctx(ctx).Info().Msg("Executing push request")

	// Find provisional image with buildID tag
	imageID, err := findImage(ctx, fmt.Sprintf("uw-provisional:%s", buildID))
	if err != nil {
		return err
	}

	// Tag provisional image with namespace/reponame:buildID

	target := fmt.Sprintf("%s/%s/%s:%s", b.registryURI, namespace, reponame, buildID)
	out, err := tagImage(ctx, imageID, target)
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			err = fmt.Errorf("failed to tag image: %s, %s", out, e.Stderr)
		}
		return fmt.Errorf("failed to tag image: %w", err)
	}

	out, err = pushImage(ctx, target)
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			err = fmt.Errorf("failed to push image: %s, %s", out, e.Stderr)
		}
		return fmt.Errorf("failed to push image: %w", err)
	}
	log.Ctx(ctx).Info().Msgf("Pushed image to %q", target)

	return nil
}

func NewBuilder(logger builder.LogDriver, registryURI string) *Builder {
	return &Builder{logger: logger, registryURI: registryURI}
}
