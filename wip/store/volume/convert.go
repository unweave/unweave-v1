package volume

import (
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
)

func dbToApi(volume db.UnweaveVolume) types.Volume {
	return types.Volume{
		ID:       volume.ID,
		Name:     volume.Name,
		Provider: types.Provider(volume.Provider),
		State: &types.VolumeState{
			CreatedAt: volume.CreatedAt,
			UpdatedAt: volume.UpdatedAt,
		},
	}
}
