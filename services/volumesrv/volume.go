package volumesrv

import (
	"context"

	"github.com/unweave/unweave/api/types"
)

type Store interface {
	VolumeAdd(projectID string, provider types.Provider, id, name string, size int) error
	VolumeList(projectID string) ([]types.Volume, error)
	VolumeGet(projectID, idOrName string) (types.Volume, error)
	VolumeDelete(id string) error
	VolumeUpdate(id string, volume types.Volume) error
}

type Driver interface {
	VolumeCreate(ctx context.Context, projectID, name string, size int) (string, error)
	VolumeDelete(ctx context.Context, id string) error
	VolumeProvider() types.Provider
	VolumeDriver(ctx context.Context) string
	VolumeResize(ctx context.Context, id string, size int) error
}

type Service interface {
	Provider() types.Provider
	Create(ctx context.Context, accountID string, projectID string, provider types.Provider, name string, size int) (types.Volume, error)
	Delete(ctx context.Context, projectID, idOrName string) error
	Get(ctx context.Context, projectID, idOrName string) (types.Volume, error)
	List(ctx context.Context, projectID string) ([]types.Volume, error)
	Resize(ctx context.Context, projectID, idOrName string, size int) error
}
