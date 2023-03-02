package server

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/unweave/unweave/db"
)

type BuilderService struct {
	srv *Service
}

func (b *BuilderService) Build(ctx context.Context, projectID string, buildCtx io.Reader) (string, error) {
	builder, err := b.srv.InitializerBuilder(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create runtime: %w", err)
	}

	buildID, err := builder.Build(ctx, buildCtx)
	if err != nil {
		return "", fmt.Errorf("failed to build image: %v", err)
	}

	params := db.BuildCreateParams{
		ID:          buildID,
		ProjectID:   projectID,
		BuilderType: builder.GetBuilder(),
		CreatedAt:   time.Time{},
	}

	if err = db.Q.BuildCreate(ctx, params); err != nil {
		return "", fmt.Errorf("failed to create build record: %v", err)
	}

	return buildID, nil
}

func (b *BuilderService) Watch(ctx context.Context, buildID string) error {

	// call builder to get build status
	// update build status in db

	return nil
}
