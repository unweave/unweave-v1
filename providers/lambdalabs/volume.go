package lambdalabs

import (
	"context"

	"github.com/unweave/unweave/api/types"
)

func (d *Driver) VolumeCreate(ctx context.Context, name string, size int) {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) VolumeDelete(ctx context.Context, id string) {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) VolumeGet(ctx context.Context, id string) (types.Volume, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) VolumeProvider() types.Provider {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) VolumeDriver(ctx context.Context) string {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) VolumeResize(ctx context.Context, id string, size int) error {
	//TODO implement me
	panic("implement me")
}
