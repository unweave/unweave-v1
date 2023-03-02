package builder

import (
	"context"
	"io"

	"github.com/unweave/unweave/api/types"
)

type Builder interface {
	GetBuilder() string
	Build(ctx context.Context, buildCtx io.Reader) (logs []types.LogEntry, err error)
}
