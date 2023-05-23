package wip

import (
	"context"
	"fmt"

	"github.com/unweave/unweave/api/server"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/vault"
	"github.com/unweave/unweave/wip/conductor"
)

// Exec is a high-level interface for executing remote code.
//
// It must be implemented by any backed that wants to support remote code execution.
type Exec interface {
	Logs() (logs chan<- types.LogEntry, err error)
}

// SSHDriver is a driver that uses SSH to manage nodes and containers running on them.
type SSHDriver struct {
	vault vault.Vault
}

func (d *SSHDriver) Create(ctx context.Context, project string, params types.ExecCreateParams) error {
	providerID, err := server.GetProviderIDFromProject(project)
	if err != nil {
		return err
	}

	spec := conductor.Spec{}
	config := conductor.ContainerCreateConfig{
		Cmd:     params.Command,
		SSHKeys: params.PublicKeys,
		Volumes: map[string]string{},
	}

	_, _, err = conductor.CreateContainer(providerID, spec, config)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	return nil
}
