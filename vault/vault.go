package vault

import "context"

type Vault interface {
	// GetSecret gets a secret from the vault.
	GetSecret(ctx context.Context, id string) (string, error)
	// DeleteSecret deletes a secret from the vault.
	DeleteSecret(ctx context.Context, id string) error
	// SetSecret sets a secret in the vault. If id is nil, a random id will be generated.
	SetSecret(ctx context.Context, secret string, id *string) (string, error)
}
