package builder

import (
	"context"
	"io"
)

type Builder interface {
	GetBuilder() string
	Build(ctx context.Context, buildCtx io.Reader) (buildID string, err error)
}
