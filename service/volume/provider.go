package volume

import (
	"context"

	"github.com/unweave/unweave/api/types"
)

type Provider interface {
	// Create creates a volume. The size is in GB.
	Create(ctx context.Context, size int) (Volume, error)
	// Delete deletes the volume.
	Delete(ctx context.Context) error
	// Get gets the volume.
	Get(ctx context.Context, id string) (Volume, error)
	// List lists all volumes the provider has.
	List(ctx context.Context) ([]Volume, error)
	// Name returns the name of the provider.
	Name() types.Provider
	// Resize resizes the volume to the given size in GB.
	Resize(ctx context.Context, size int) error
}

// Volume is an interface that must be implemented by a volume.
// The implementation of the volume interface should have the ability/permissions
// to mount the volume onto a container inside a node. For instance, this could
// be an AWS EBS volume mounted onto a container inside an EC2 machine, a FUSE
// volume on top of S3/GCS, etc.
type Volume interface {
	// ID returns the ID of the volume.
	ID() string
	// Mount mounts the volume to the given path.
	Mount(ctx context.Context, path string) error
	// Unmount unmounts the volume from the given path.
	Unmount(ctx context.Context, path string) error
}
