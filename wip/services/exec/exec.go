package exec

import (
	"context"

	"github.com/unweave/unweave/api/types"
)

// Stats represents the resource usage of an exec.
type Stats struct {
	CPU  float64
	Mem  float64
	Disk float64
	GPU  float64
}

type Store interface {
	Create(project string, exec types.Exec) error
	Get(id string) (types.Exec, error)
	GetDriver(id string) (string, error)
	List(project string) ([]types.Exec, error)
	ListAll() ([]types.Exec, error)
	Delete(project, id string) error
	Update(id string, exec types.Exec) error
}

type Driver interface {
	Create(ctx context.Context, project, image string, spec types.HardwareSpec) (string, error)
	DriverName() string
	Get(ctx context.Context, id string) (types.Exec, error)
	List(ctx context.Context, project string) ([]types.Exec, error)
	Terminate(ctx context.Context, id string) error
	Stats(ctx context.Context, id string) (Stats, error)
}
