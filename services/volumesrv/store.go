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

func (p postgresStore) VolumeAdd(projectID string, volume types.Volume) error {
	ctx := context.Background()
	_, err := db.Q.VolumeCreate(ctx, db.VolumeCreateParams{
		ID:        volume.ID,
		ProjectID: projectID,
		Provider:  volume.Provider.String(),
	})
	if err != nil {
		return &types.Error{
			Code:    http.StatusBadRequest,
			Message: "Volume could not be created",
			Err:     err,
		}
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

	out := make([]types.Volume, 0, len(vols))
	for _, v := range vols {
		out = append(out, volumeFromDB(v))
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
				Code:    http.StatusBadRequest,
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
		return &types.Error{
			Code:    http.StatusBadRequest,
			Message: "Volume could not be deleted",
			Err:     err,
		}
	}

	return nil
}

func (p postgresStore) VolumeUpdate(id string, volume types.Volume) error {
	ctx := context.Background()
	err := db.Q.VolumeUpdate(ctx, db.VolumeUpdateParams{
		ID:   id,
		Size: int32(volume.Size),
	})
	if err != nil {
		return &types.Error{
			Code:    http.StatusBadRequest,
			Message: "Volume could not be updated",
			Err:     err,
		}
	}

	return nil
}
