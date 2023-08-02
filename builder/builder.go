package builder

import (
	"context"
	"io"

	"github.com/unweave/unweave-v1/api/types"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o builderfakes . LogDriver
//counterfeiter:generate -o builderfakes . Builder

// LogDriver defines the interface for storing and retrieving build logs.
type LogDriver interface {
	// GetLogs returns the logs for a build.
	GetLogs(ctx context.Context, buildID string) (logs []types.LogEntry, err error)
	// SaveLogs saves the logs for a build in long term storage.
	SaveLogs(ctx context.Context, buildID string, logs []types.LogEntry) error
}

// Builder defines the interface for building and storing container images.
type Builder interface {
	// BuildAndPush builds a container image from a build context and pushes it
	// to the cointainer registry.  The build context is a zip file containing
	// the source code and any other files needed to build the image.
	BuildAndPush(ctx context.Context, buildID, namespace, reponame string, buildCtx io.Reader) error
	// GetBuilder returns the name of the builder.
	GetBuilder() string
	// GetImageURI returns the URI of the image in the container registry.
	GetImageURI(ctx context.Context, buildID, namespace, reponame string) string
	// Logs returns the logs for a build.
	Logs(ctx context.Context, buildID string) (logs []types.LogEntry, err error)
}
