package conductor

import (
	"github.com/unweave/unweave-v1/wip/conductor/node"
	"github.com/unweave/unweave-v1/wip/conductor/volume"
)

type Keys interface {
	KeysRegister()
	KeysUnregister()
	KeysList()
}

type Network interface {
	NetworkAddPort()
	NetworkRemovePort()
}

type Provider interface {
	// ID is the unique identifier for the provider. This is used to monitor, assign and
	// forward requests to each provider. You should make sure this is unique. In most
	// cases this will be the accountID or userID of the owner of the provider.
	ID() string
	// Name is the name of the provider (eg. AWS, GCP, DigitalOcean etc.)
	Name() string
	node.Provider
	volume.Provider
}
