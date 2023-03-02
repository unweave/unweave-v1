package server

import (
	"context"
	"fmt"
	"io"
)

type BuilderService struct {
	srv *Service
}

func (b *BuilderService) Build(ctx context.Context, buildCtx io.Reader) (string, error) {
	builder, err := b.srv.InitializerBuilder(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create runtime: %w", err)
	}

	buildID, err := builder.Build(ctx, buildCtx)
	if err != nil {
		return "", fmt.Errorf("failed to build image: %v", err)
	}

	// save buildID to db

	return buildID, nil
}
