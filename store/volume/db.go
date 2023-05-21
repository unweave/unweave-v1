package volume

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
)

type PostgresStore struct{}

func (v *PostgresStore) Create(namespace string, id, provider string) (types.Volume, error) {
	ctx := context.Background()
	params := db.VolumeCreateParams{
		ID:        id,
		ProjectID: namespace,
		Provider:  provider,
	}
	vol, err := db.Q.VolumeCreate(ctx, params)
	if err != nil {
		return types.Volume{}, fmt.Errorf("failed to create volume in db: %w", err)
	}

	return dbToApi(vol), nil
}

func (v *PostgresStore) Get(namespace, id string) (types.Volume, error) {
	ctx := context.Background()
	vol, err := db.Q.VolumeGet(ctx, id)
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

	return dbToApi(vol), nil

}

func (v *PostgresStore) List(namespace string) ([]types.Volume, error) {
	ctx := context.Background()
	vols, err := db.Q.VolumeList(ctx, namespace)
	if err != nil {
		if err == sql.ErrNoRows {
			return []types.Volume{}, nil
		}
		return []types.Volume{}, fmt.Errorf("failed to list volumes from db: %w", err)
	}

	apiVols := make([]types.Volume, len(vols))
	for i, vol := range vols {
		apiVols[i] = dbToApi(vol)
	}
	return apiVols, nil
}

func (v *PostgresStore) Remove(namespace, id string) error {
	return nil
}

func (v *PostgresStore) Update(namespace string, volume types.Volume) error {
	return nil
}
