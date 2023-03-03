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
		return nil, nil, fmt.Errorf("buildPath %q does not exist: %v", buildPath, err)
	}

	c := []string{
		"docker",
		"build",
		//"--cache-from", cache,
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

	// Buffer channel in case i/o is slow
	logsch = make(chan string, 1000)
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
		}
		if e := cmd.Wait(); e != nil {
			errch <- e
		}
	}()

	return logsch, errch, nil
}

// pushImage pushes the image to the registry
func pushImage(ctx context.Context, image string) (output string, err error) {
	cmd := exec.CommandContext(
		ctx,
		"docker",
		"push",
		image,
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

type Builder struct{}

func (b *Builder) GetBuilder() string {
	return "docker"
}

func (b *Builder) Build(ctx context.Context, buildID string, buildCtx io.Reader) ([]types.LogEntry, error) {
	ctx = log.With().
		Str("builder", b.GetBuilder()).
		Str("buildID", buildID).
		Logger().WithContext(ctx)

	log.Ctx(ctx).Info().Msg("Executing build request")

	dir := fmt.Sprintf("/tmp/uw-build-ctx/%s", buildID)
	buildBytes, err := io.ReadAll(buildCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to read build context: %v", err)
	}

	if err := saveContext(dir, buildBytes); err != nil {
		return nil, fmt.Errorf("failed to save build context: %v", err)
	}

	imageURI := fmt.Sprintf("uw-provisional:%s", buildID) // until tagged
	logsch, errch, err := buildImage(ctx, dir, imageURI, "")
	if err != nil {
		return nil, fmt.Errorf("failed to build image: %v", err)
	}
	log.Ctx(ctx).Info().Msg("Started image build with build context at " + dir)

	var logs []types.LogEntry

	for {
		select {
		case <-ctx.Done():
			return logs, nil
		case l, ok := <-logsch:
			if !ok {
				return logs, nil
			}
			logs = append(logs, types.LogEntry{TimStamp: time.Now(), Message: l})
		case e := <-errch:
			log.Ctx(ctx).Error().Err(e).Msg("Error building image")
		}
	}
}

func (b *Builder) Push(ctx context.Context, repo, tag string) error {
	//
	return nil
}
