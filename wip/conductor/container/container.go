package container

import (
	"context"
	"time"

	"github.com/unweave/unweave/wip/conductor/types"
)

type Container struct {
	ID       string
	ImageSHA string
	Cmd      []string
	State    types.ContainerState
	Network  types.Network
	Spec     types.Spec
}

type LogEntry struct {
	ContainerID string    `json:"containerID"`
	NodeID      string    `json:"nodeID"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
	Level       string    `json:"level"`
}

// Driver is an interface that must be implemented by any type that needs access to
// high level container operations.
// The implementation should handle the details of operating over a network, different
// cloud providers etc.
type Driver interface {
	Logs(ctx context.Context, id string) (logs chan<- LogEntry, err error)
	Exec(ctx context.Context, id string, cmd []string) error
	Stop(ctx context.Context, id string) error
	Kill(ctx context.Context, id string) error
}

// Runtime is an interface that must be implemented by any runtime that manages containers
// on a node. Docker is a classic example of this.
type Runtime interface {
	ContainerCreate(ctx context.Context) (string, error)
	ContainerExec(ctx context.Context, containerID string)
	ContainerLogs(ctx context.Context, containerID string)
	ContainerList(ctx context.Context) ([]Container, error)
	ContainerStart(ctx context.Context, containerID string)
	ContainerStop(ctx context.Context, containerID string) error
	ContainerTop(ctx context.Context, containerID string)
	ContainerKill(ctx context.Context, containerID string)
}
