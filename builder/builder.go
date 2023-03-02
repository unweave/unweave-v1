package builder

import (
	"context"
	"io"
)

type Builder interface {
	Build(ctx context.Context, buildCtx io.Reader) (buildID string, err error)
}
