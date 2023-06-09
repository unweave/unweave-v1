package volume

import "github.com/unweave/unweave/api/types"

// Store is an interface that must be implemented any volume store.
type Store interface {
	Create(namespace string, id, provider string) (types.Volume, error)
	Get(namespace, id string) (types.Volume, error)
	List(namespace string) ([]types.Volume, error)
	Remove(namespace, id string) error
	Update(namespace string, volume types.Volume) error
}
