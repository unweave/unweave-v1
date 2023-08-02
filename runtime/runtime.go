// Package runtime manages the lifecycle of a session.
package runtime

import (
	"context"

	"github.com/unweave/unweave-v1/builder"
	"github.com/unweave/unweave-v1/vault"
)

type Initializer interface {
	InitializeBuilder(ctx context.Context, userID string, builder string) (builder.Builder, error)
	InitializeVault(ctx context.Context) (vault.Vault, error)
}
