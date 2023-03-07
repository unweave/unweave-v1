package builder

import (
	"context"
	"io"

	"github.com/unweave/unweave/api/types"
)

type LogDriver interface {
	// GetLogs returns the logs for a build.
	GetLogs(ctx context.Context, buildID string) (logs []types.LogEntry, err error)
	// SaveLogs saves the logs for a build in long term storage.
	SaveLogs(ctx context.Context, buildID string, logs []types.LogEntry) error
}

type Builder interface {
	GetBuilder() string
	// Build builds a container image from a build context.
	// The build context is a zip file containing the source code and any other files
	// needed to build the image.
	Build(ctx context.Context, buildID string, buildCtx io.Reader) (err error)
	// Logs returns the logs for a build.
	Logs(ctx context.Context, buildID string) (logs []types.LogEntry, err error)
	// Push pushes an image to the container registry. The buildID is used as the tag.
	// If you want to use a different tag, use the Tag method instead.
	Push(ctx context.Context, buildID, namespace, reponame string) error
}
