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
//
// The implementing type should handle the details of operating over a network, different
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
	Create(ctx context.Context) (string, error)
	Exec(ctx context.Context, containerID string)
	Logs(ctx context.Context, containerID string)
	List(ctx context.Context) ([]Container, error)
	Start(ctx context.Context, containerID string)
	Stop(ctx context.Context, containerID string) error
	Top(ctx context.Context, containerID string)
	Kill(ctx context.Context, containerID string)
}
