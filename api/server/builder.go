package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
)

type BuildMetaDataV1 struct {
	Version string `json:"version"`
}

type BuilderService struct {
	srv *Service
}

func (b *BuilderService) Build(ctx context.Context, projectID string, params *types.BuildsCreateParams) (string, error) {
	builder, err := b.srv.InitializeBuilder(ctx, params.Builder)
	if err != nil {
		return "", fmt.Errorf("failed to create runtime: %w", err)
	}

	bcp := db.BuildCreateParams{
		ProjectID:   projectID,
		BuilderType: builder.GetBuilder(),
	}

	buildID, err := db.Q.BuildCreate(ctx, bcp)
	if err != nil {
		return "", fmt.Errorf("failed to create build record: %v", err)
	}

	go func() {
		c := context.Background()
		c = log.With().Str(BuildIDCtxKey, buildID).Logger().WithContext(c)

		err := builder.Build(c, buildID, params.BuildContext)
		if err != nil {
			p := db.BuildUpdateParams{ID: buildID}

			var e *types.Error
			if errors.As(err, &e) && e.Code == http.StatusBadRequest {
				log.Ctx(c).Warn().Err(err).Msg("User build failed")
				p.Status = db.UnweaveBuildStatusFailed
			} else {
				log.Ctx(c).Error().Err(err).Msg("Failed to build image")
				p.Status = db.UnweaveBuildStatusError
			}

			meta, err := json.Marshal(BuildMetaDataV1{Version: "1"})
			if err != nil {
				log.Ctx(c).Error().Err(err).Msg("Failed to marshal build metadata")
			}
			p.MetaData = meta

			if err := db.Q.BuildUpdate(c, p); err != nil {
				log.Ctx(c).Error().Err(err).Msg("Failed to set build error in DB")
			}
			return
		}

		meta, err := json.Marshal(BuildMetaDataV1{Version: "1"})
		if err != nil {
			log.Ctx(c).Error().Err(err).Msg("Failed to marshal build metadata")
		}

		p := db.BuildUpdateParams{
			ID:       buildID,
			Status:   db.UnweaveBuildStatusSuccess,
			MetaData: meta,
		}
		if err := db.Q.BuildUpdate(c, p); err != nil {
			log.Ctx(c).Error().Err(err).Msg("Failed to set build success in DB")
		}
	}()

	return buildID, nil
}

func (b *BuilderService) GetLogs(ctx context.Context, buildID string) ([]types.LogEntry, error) {
	build, err := db.Q.BuildGet(ctx, buildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get build: %v", err)
	}

	builder, err := b.srv.InitializeBuilder(ctx, build.BuilderType)
	if err != nil {
		return nil, fmt.Errorf("failed to initializer builder: %w", err)
	}

	logs, err := builder.Logs(ctx, buildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs from builder: %w", err)
	}
	return logs, nil
}

func (b *BuilderService) Watch(ctx context.Context, buildID string) error {

	// call builder to get build status
	// update build status in db

	return nil
}
