package volumesrv

import (
	"github.com/unweave/unweave/api/types"
)

type postgresStore struct{}

func NewPostgresStore() Store {
	return postgresStore{}
}

func (p postgresStore) VolumeAdd(projectID string, volume types.Volume) error {
	//TODO implement me
	panic("implement me")
}

func (p postgresStore) VolumeGet(projectID, idOrName string) (types.Volume, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresStore) VolumeDelete(id string) {
	//TODO implement me
	panic("implement me")
}

func (p postgresStore) VolumeUpdate(id string, volume types.Volume) error {
	//TODO implement me
	panic("implement me")
}
