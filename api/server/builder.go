package server

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/tools/random"
)

func handleBuildErr(ctx context.Context, buildID string, err error) {
	p := db.BuildUpdateParams{
		ID:     buildID,
		Status: "building",
	}

	var e *types.Error
	var errmeta string
	if errors.As(err, &e) && e.Code == http.StatusBadRequest {
		log.Ctx(ctx).Warn().Err(err).Msg("User build failed")
		p.Status = db.UnweaveBuildStatusFailed
		errmeta = fmt.Sprintf("Build failed: %v", e.Message)
	} else {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to build image")
		p.Status = db.UnweaveBuildStatusError
		errmeta = fmt.Sprintf("Build error: Something went wrong. Please contact us for support.")
	}
	meta, merr := json.Marshal(BuildMetaDataV1{
		Version: 1,
		Error:   errmeta,
	})
	if merr != nil {
		log.Ctx(ctx).Error().Err(merr).Msg("Failed to marshal build metadata")
	}
	p.MetaData = meta

	if derr := db.Q.BuildUpdate(ctx, p); derr != nil {
		log.Ctx(ctx).Error().Err(derr).Msg("Failed to set build error in DB")
	}
}

func convertZipToTarGz(zipReader io.Reader) (io.Reader, error) {
	zipData, err := io.ReadAll(zipReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read zip file: %w", err)
	}

	zipFileReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %w", err)
	}

	foundDockerfile := false

	tarGzBuffer := &bytes.Buffer{}
	gzipWriter := gzip.NewWriter(tarGzBuffer)
	tarWriter := tar.NewWriter(gzipWriter)

	defer gzipWriter.Close()
	defer tarWriter.Close()

	for _, zipFile := range zipFileReader.File {
		fileReader, err := zipFile.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file from zip: %w", err)
		}

		header := &tar.Header{
			Name:     filepath.ToSlash(zipFile.Name),
			Mode:     int64(zipFile.Mode()),
			Size:     int64(zipFile.UncompressedSize64),
			ModTime:  zipFile.Modified,
			Typeflag: tar.TypeReg,
		}
		if zipFile.FileInfo().IsDir() {
			header.Typeflag = tar.TypeDir
			header.Name = strings.TrimSuffix(header.Name, "/") + "/"
		}

		if header.Name == "Dockerfile" {
			foundDockerfile = true
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return nil, fmt.Errorf("failed to write tar header: %w", err)
		}
		if _, err := io.Copy(tarWriter, fileReader); err != nil {
			return nil, fmt.Errorf("failed to copy file contents to tar: %w", err)
		}

		fileReader.Close()
	}

	if !foundDockerfile {
		return nil, &types.Error{
			Code:       http.StatusBadRequest,
			Message:    "No Dockerfile found in build context",
			Suggestion: "Make sure your build context contains a Dockerfile",
			Err:        fmt.Errorf("no Dockerfile found in build context"),
		}
	}

	return tarGzBuffer, nil
}

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

	buildContext, err := convertZipToTarGz(params.BuildContext)
	if err != nil {
		return "", fmt.Errorf("failed to convert build context: %w", err)
	}

	buildID, err := db.Q.BuildCreate(ctx, bcp)
	if err != nil {
		return "", fmt.Errorf("failed to create build record: %v", err)
	}

	go func() {
		c := context.Background()
		c = log.With().Str(types.BuildIDCtxKey, buildID).Logger().WithContext(c)

		// Upload context to S3
		if e := builder.Upload(c, buildID, buildContext); e != nil {
			log.Ctx(c).Error().Err(e).Msgf("Failed to upload build context for build %q", buildID)
			handleBuildErr(c, buildID, e)
			return
		}

		// Reponame must be lowercase for dockerhub
		reponame := strings.ToLower(projectID)
		namespace := strings.ToLower(b.srv.aid)

		if e := builder.BuildAndPush(c, buildID, namespace, reponame, params.BuildContext); e != nil {
			handleBuildErr(c, buildID, e)
			return
		}

		meta, err := json.Marshal(BuildMetaDataV1{Version: 1})
		if err != nil {
			log.Ctx(c).Error().Err(err).Msg("Failed to marshal build metadata")
		}

		p := db.BuildUpdateParams{
			ID:         buildID,
			Status:     db.UnweaveBuildStatusSuccess,
			FinishedAt: time.Now(),
			MetaData:   meta,
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

func (b *BuilderService) GetImageURI(ctx context.Context, buildID string) (string, error) {
	build, err := db.Q.BuildGet(ctx, buildID)
	if err != nil {
		return "", fmt.Errorf("failed to get build: %v", err)
	}

	builder, err := b.srv.InitializeBuilder(ctx, build.BuilderType)
	if err != nil {
		return "", fmt.Errorf("failed to initializer builder: %w", err)
	}

	reponame := strings.ToLower(build.ProjectID) // reponame must be lowercase for dockerhub
	namespace := strings.ToLower(b.srv.aid)
	uri := builder.GetImageURI(ctx, build.ID, namespace, reponame)

	return uri, nil
}
