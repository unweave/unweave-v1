package local

import (
	"context"

	"github.com/unweave/unweave-v1/wip/conductor/node"
	"github.com/unweave/unweave-v1/wip/conductor/volume"
)

// Local provider runs locally on the machine. This should only be used for testing.
type Local struct {
	id     string
	driver *NodeDriver
}

func NewProvider(id string) *Local {
	return &Local{
		driver: &NodeDriver{},
	}
}

func (l *Local) ID() string {
	return l.id
}

func (l *Local) Name() string {
	return "local"
}

func (l *Local) NodeCreate(ctx context.Context) (string, error) {
	return "", nil
}

func (l *Local) NodeDelete(ctx context.Context, id string) error {
	return nil
}

func (l *Local) NodeInit(ctx context.Context, id string, options ...func(node.Driver)) (*node.Node, error) {
	return nil, nil
}

func (l *Local) NodeList(ctx context.Context) ([]node.Node, error) {
	return nil, nil
}

func (l *Local) VolumeCreate(ctx context.Context, size int) (volume.Volume, error) {
	return nil, nil
}

func (l *Local) VolumeDelete(ctx context.Context) error {
	return nil
}

func (l *Local) VolumeGet(ctx context.Context, id string) (volume.Volume, error) {
	return nil, nil
}

func (l *Local) VolumeList(ctx context.Context) ([]volume.Volume, error) {
	return nil, nil
}

func (l *Local) VolumeResize(ctx context.Context, size int) error {
	return nil
}
