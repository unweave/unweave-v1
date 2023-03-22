package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/tools/random"
)

// BuildMetaDataV1 versions the metadata for a build stored in the DB.
type BuildMetaDataV1 struct {
	Version int16  `json:"version"`
	Error   string `json:"error"`
}

type BuilderService struct {
	srv *Service
}

func (b *BuilderService) Build(ctx context.Context, projectID string, params *types.BuildsCreateParams) (string, error) {
	builder, err := b.srv.InitializeBuilder(ctx, params.Builder)
	if err != nil {
		return "", fmt.Errorf("failed to create runtime: %w", err)
	}

	if params.Name == nil {
		name := random.GenerateRandomAdjectiveNounTriplet()
		params.Name = &name
	}

	bcp := db.BuildCreateParams{
		ProjectID:   projectID,
		BuilderType: builder.GetBuilder(),
		Name:        *params.Name,
		CreatedBy:   b.srv.cid,
	}

	buildID, err := db.Q.BuildCreate(ctx, bcp)
	if err != nil {
		return "", fmt.Errorf("failed to create build record: %v", err)
	}

	go func() {
		c := context.Background()
		c = log.With().Str(BuildIDCtxKey, buildID).Logger().WithContext(c)

		// Build

		err := builder.Build(c, buildID, params.BuildContext)
		if err != nil {
			log.Ctx(c).Error().Err(err).Msg("Failed to build image")

			p := db.BuildUpdateParams{ID: buildID}

			var e *types.Error
			var errmeta string
			if errors.As(err, &e) && e.Code == http.StatusBadRequest {
				log.Ctx(c).Warn().Err(err).Msg("User build failed")
				p.Status = db.UnweaveBuildStatusFailed
				errmeta = fmt.Sprintf("Build failed: %v", err.Error())
			} else {
				log.Ctx(c).Error().Err(err).Msg("Failed to build image")
				p.Status = db.UnweaveBuildStatusError
				errmeta = fmt.Sprintf("Build error: %v", err.Error())
			}

			meta, merr := json.Marshal(BuildMetaDataV1{
				Version: 1,
				Error:   errmeta,
			})
			if merr != nil {
				log.Ctx(c).Error().Err(merr).Msg("Failed to marshal build metadata")
			}
			p.MetaData = meta

			if derr := db.Q.BuildUpdate(c, p); derr != nil {
				log.Ctx(c).Error().Err(derr).Msg("Failed to set build error in DB")
			}
			return
		}

		// Push

		reponame := strings.ToLower(projectID) // reponame must be lowercase for dockerhub
		namespace := strings.ToLower(b.srv.cid)
		err = builder.Push(c, buildID, namespace, reponame)
		if err != nil {
			log.Ctx(c).Error().Err(err).Msg("Failed to push image")

			meta, e := json.Marshal(BuildMetaDataV1{
				Version: 1,
				Error:   fmt.Sprintf("Builc push failed: %v", err.Error()),
			})
			if e != nil {
				log.Ctx(c).Error().Err(e).Msg("Failed to marshal build metadata")
			}
			p := db.BuildUpdateParams{
				ID:       buildID,
				Status:   db.UnweaveBuildStatusError,
				MetaData: meta,
			}
			if e := db.Q.BuildUpdate(c, p); e != nil {
				log.Ctx(c).Error().Err(e).Msg("Failed to set build error in DB")
			}
			return
		}

		meta, err := json.Marshal(BuildMetaDataV1{Version: 1})
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
