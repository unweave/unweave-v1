package volumesrv

import (
	"github.com/unweave/unweave-v1/api/types"
	"github.com/unweave/unweave-v1/db"
)

func volumeFromDB(volume db.UnweaveVolume) types.Volume {
	return types.Volume{
		ID:   volume.ID,
		Name: volume.Name,
		Size: int(volume.Size),
		State: types.VolumeState{
			CreatedAt: volume.CreatedAt,
			UpdatedAt: volume.UpdatedAt,
		},
		Provider: types.Provider(volume.Provider),
	}
}
