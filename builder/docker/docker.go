package docker

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
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

	cmd := exec.CommandContext(
		ctx,
		"docker",
		"build",
		"--cache-from", cache,
		"--build-arg", "BUILDKIT_INLINE_CACHE=1",
		"-t", image,
		buildPath,
	)
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

	func() {
		if e := cmd.Start(); e != nil {
			errch <- e
		}
		if e := cmd.Wait(); e != nil {
			errch <- e
		}
	}()

	return logsch, errch, nil
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

type Builder struct{}

func (b *Builder) Build(ctx context.Context, buildCtx io.Reader) (string, error) {
	log.Ctx(ctx).Info().Msg("Building docker image")

	// write build context to disk
	// build image

	logsch, errch, err := buildImage(ctx, "", "", "")
	if err != nil {
		return "", fmt.Errorf("failed to build image: %v", err)
	}

	go func() {
		for {
			select {
			case l := <-logsch:
				log.Ctx(ctx).Info().Msg(l)
			case e := <-errch:
				log.Ctx(ctx).Error().Msg(e.Error())
			}
		}
	}()

	buildID := "bld_123"
	return buildID, nil
}

func (b *Builder) Push(ctx context.Context, repo, tag string) error {
	//
	return nil
}
