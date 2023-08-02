//nolint:revive
package lambdalabs

import (
	"context"

	"github.com/unweave/unweave-v1/api/types"
)

func (d *Driver) VolumeCreate(ctx context.Context, projectID, name string, size int) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Driver) VolumeDelete(ctx context.Context, id string) error {
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
