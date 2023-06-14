package volumesrv

import (
	"context"

	"github.com/unweave/unweave/api/types"
)

type Store interface {
	VolumeAdd(projectID string, volume types.Volume) error
	VolumeList(projectID string) ([]types.Volume, error)
	VolumeGet(projectID, idOrName string) (types.Volume, error)
	VolumeDelete(id string) error
	VolumeUpdate(id string, volume types.Volume) error
}

type Driver interface {
	VolumeCreate(ctx context.Context, name string, size int) error
	VolumeDelete(ctx context.Context, id string) error
	VolumeGet(ctx context.Context, id string) (types.Volume, error)
	VolumeProvider() types.Provider
	VolumeDriver(ctx context.Context) string
	VolumeResize(ctx context.Context, id string, size int) error
}
