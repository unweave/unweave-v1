package exec

import (
	"context"

	"github.com/unweave/unweave/api/types"
)

type Store interface {
	Create(project string, exec types.Exec) error
	Get(id string) (types.Exec, error)
	GetDriver(id string) (string, error)
	List(project string) ([]types.Exec, error)
	ListByProvider(provider types.Provider, filterActive bool) ([]types.Exec, error)
	Delete(project, id string) error
	Update(id string, exec types.Exec) error
}

type Driver interface {
	Create(ctx context.Context, project, image string, spec types.HardwareSpec, pubKeys []string, region *string) (string, error)
	DriverName() string
	Provider() types.Provider
	Terminate(ctx context.Context, id string) error
	Spec(ctx context.Context, id string) (types.HardwareSpec, error)
	Stats(ctx context.Context, id string) (Stats, error)
}
