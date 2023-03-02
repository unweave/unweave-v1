package builder

import (
	"context"
	"io"
)

type Builder interface {
	GetBuildLogs(ctx context.Context, imageID string) (io.ReadCloser, error)
	Build(ctx context.Context, buildCtx io.Reader) (buildID string, err error)
	Push(ctx context.Context, repo, tag string) error
}
