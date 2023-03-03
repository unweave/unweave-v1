package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
)

type BuildMetaDataV1 struct {
	Version string           `json:"version"`
	Logs    []types.LogEntry `json:"logs"`
}

type BuilderService struct {
	srv *Service
}

func (b *BuilderService) Build(ctx context.Context, projectID string, params *types.ImageBuildParams) (string, error) {
	builder, err := b.srv.InitializerBuilder(ctx, params.Builder)
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

		logs, e := builder.Build(c, buildID, params.BuildContext)
		if e != nil {
			log.Ctx(c).Error().Err(e).Msg("Failed to build image")
			meta, err := json.Marshal(BuildMetaDataV1{Version: "1", Logs: logs})
			if err != nil {
				log.Ctx(c).Error().Err(err).Msg("Failed to marshal build metadata")
			}

			p := db.BuildUpdateParams{
				ID:       buildID,
				Status:   db.UnweaveBuildStatusError,
				MetaData: meta,
			}
			if err := db.Q.BuildUpdate(c, p); err != nil {
				log.Ctx(c).Error().Err(err).Msg("Failed to set build error in DB")
			}
			return
		}

		meta, err := json.Marshal(BuildMetaDataV1{Version: "1", Logs: logs})
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

func (b *BuilderService) GetLogs(ctx context.Context, buildID string) (io.ReadCloser, error) {
	return nil, nil
}

func (b *BuilderService) Watch(ctx context.Context, buildID string) error {

	// call builder to get build status
	// update build status in db

	return nil
}
