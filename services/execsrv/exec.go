package execsrv

import (
	"context"
	"errors"

	"github.com/unweave/unweave/api/types"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Store interface {
	Create(project string, exec types.Exec) error
	Get(id string) (types.Exec, error)
	GetDriver(id string) (string, error)
	List(filterProject *string, filterProvider *types.Provider, filterActive bool) ([]types.Exec, error)
	Delete(id string) error
	Update(id string, exec types.Exec) error
	UpdateStatus(id string, status types.Status) error
}

type Driver interface {
	ExecCreate(ctx context.Context, project, image string, spec types.HardwareSpec, network types.NetworkSpec, volumes []types.ExecVolume, pubKeys []string, region *string) (string, error)
	ExecDriverName() string
	ExecGetStatus(ctx context.Context, execID string) (types.Status, error)
	ExecProvider() types.Provider
	ExecTerminate(ctx context.Context, id string) error
	ExecSpec(ctx context.Context, id string) (types.HardwareSpec, error)
	ExecStats(ctx context.Context, id string) (Stats, error)
	// ExecPing pings the driver availability on behalf of a user. This can be used to
	// check if the driver is configured correctly and healthy.
	ExecPing(ctx context.Context, accountID *string) error
}
