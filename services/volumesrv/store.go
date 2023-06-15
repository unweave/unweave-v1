package volumesrv

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
)

type postgresStore struct{}

func NewPostgresStore() Store {
	return postgresStore{}
}

func (p postgresStore) VolumeAdd(projectID string, provider types.Provider, id string, name string, size int) error {
	ctx := context.Background()
	params := db.VolumeCreateParams{
		ID:        id,
		ProjectID: projectID,
		Provider:  provider.String(),
		Name:      name,
		Size:      int32(size),
	}
	_, err := db.Q.VolumeCreate(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create volume in db: %w", err)
	}

	return nil
}

func (p postgresStore) VolumeList(projectID string) ([]types.Volume, error) {
	vols, err := db.Q.VolumeList(context.Background(), projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get volumes from db: %w", err)
	}

	out := make([]types.Volume, len(vols))
	for idx, v := range vols {
		out[idx] = volumeFromDB(v)
	}

	return out, nil
}

func (p postgresStore) VolumeGet(projectID, idOrName string) (types.Volume, error) {
	vol, err := db.Q.VolumeGet(context.Background(), db.VolumeGetParams{
		ProjectID: projectID,
		ID:        idOrName,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Volume{}, &types.Error{
				Code:    http.StatusNotFound,
				Message: "Volume not found",
				Err:     err,
			}
		}
		return types.Volume{}, fmt.Errorf("failed to get volume from db: %w", err)
	}

	return volumeFromDB(vol), nil
}

func (p postgresStore) VolumeDelete(id string) error {
	ctx := context.Background()
	if err := db.Q.VolumeDelete(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			return &types.Error{
				Code:    http.StatusNotFound,
				Message: "Volume not found",
				Err:     err,
			}
		}
		return fmt.Errorf("failed to delete volume from db: %w", err)
	}

	return nil
}

func (p postgresStore) VolumeUpdate(id string, volume types.Volume) error {
	ctx := context.Background()
	params := db.VolumeUpdateParams{
		ID:   id,
		Size: int32(volume.Size),
	}

	err := db.Q.VolumeUpdate(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return &types.Error{
				Code:    http.StatusNotFound,
				Message: "Volume not found",
				Err:     err,
			}
		}
		return fmt.Errorf("failed to update volume in db: %w", err)
	}

	return nil
}
