package runtime

import "context"

type Vault interface {
	// Encrypt encrypts data and stores it under the provided key. It can only be
	// decrypted using the Decrypt function with the same key. The returned cipher
	// is base64 encoded.
	Encrypt(ctx context.Context, key string, data []byte) (cipher string, err error)
	// Decrypt decrypts data encrypted with the Encrypt function. The cipher must be
	// base64 encoded.
	Decrypt(ctx context.Context, key string, cipher string) (data []byte, err error)
	// GetSecret returns the secret value for the given key.
	GetSecret(ctx context.Context, key string) (string, error)
	// SetSecret sets the secret value for the given key.
	SetSecret(ctx context.Context, key string, value string) error
}
