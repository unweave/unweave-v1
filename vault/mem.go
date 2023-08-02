package vault

import (
	"context"
	"fmt"
	"sync"

	"github.com/unweave/unweave-v1/tools"
	"github.com/unweave/unweave-v1/tools/random"
)

// MemVault is an in-memory implementation of the Vault interface.
type MemVault struct {
	store map[string]string
	mutex *sync.Mutex
}

func (m *MemVault) GetSecret(ctx context.Context, id string) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.store[id], nil
}

func (m *MemVault) DeleteSecret(ctx context.Context, id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.store, id)
	return nil
}

func (m *MemVault) SetSecret(ctx context.Context, secret string, id *string) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if id == nil {
		id = tools.Stringy("scr_" + random.GenerateRandomPhrase(10, "-"))
	}
	if _, ok := m.store[*id]; ok {
		return "", fmt.Errorf("secret with id %s already exists", *id)
	}
	m.store[*id] = secret
	return *id, nil
}

func NewMemVault() *MemVault {
	return &MemVault{
		store: map[string]string{},
		mutex: &sync.Mutex{},
	}
}
