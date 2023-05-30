package lambdalabs

import (
	"context"

	"github.com/unweave/unweave/api/types"
	execsrv "github.com/unweave/unweave/wip/services/exec"
)

// ExecDriver implements the exec.Driver interface for Lambda Labs.
// Lambda Labs needs a special implementation since they don't support Docker on their
// VMs. This means we currently can't run an Exec as a container and instead default to
// the bare VM with the pre-configured Lambda Labs image.
type ExecDriver struct{}

func (e ExecDriver) Create(ctx context.Context, project, image string, spec types.HardwareSpec, pubKeys []string) (string, error) {
	//TODO implement me
	// TODO: we need to check if the public key already exists on the VM and if not, add it
	panic("implement me")
}

func (e ExecDriver) DriverName() string {
	//TODO implement me
	panic("implement me")
}

func (e ExecDriver) Get(ctx context.Context, id string) (types.Exec, error) {
	//TODO implement me
	panic("implement me")
}

func (e ExecDriver) List(ctx context.Context, project string) ([]types.Exec, error) {
	//TODO implement me
	panic("implement me")
}

func (e ExecDriver) Provider() types.Provider {
	return types.LambdaLabsProvider
}

func (e ExecDriver) Terminate(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (e ExecDriver) Stats(ctx context.Context, id string) (execsrv.Stats, error) {
	//TODO implement me
	panic("implement me")
}
