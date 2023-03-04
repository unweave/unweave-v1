package builder

import (
	"context"
	"io"

	"github.com/unweave/unweave/api/types"
)

type Builder interface {
	GetBuilder() string
	// GetLogs returns the logs for a build.
	GetLogs(ctx context.Context, buildID string) (logs []types.LogEntry, err error)
	// Build builds a container image from a build context.
	// The build context is a zip file containing the source code and any other files
	// needed to build the image.
	Build(ctx context.Context, buildID string, buildCtx io.Reader) (logs []types.LogEntry, err error)
	// SaveLogs saves the logs for a build in long term storage.
	SaveLogs(ctx context.Context, buildID string, logs []types.LogEntry) error
}
