package ssh_keys

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/unweave/unweave/db"
)

type Store struct{}

func (s Store) GetSSHKeyByNameIfExists(ctx context.Context, name, userID string) (*db.UnweaveSshKey, error) {
	p := db.SSHKeyGetByNameParams{Name: name, OwnerID: userID}
	key, err := db.Q.SSHKeyGetByName(ctx, p)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get SSH key from DB: %w", err)
	}

	return &key, err
}

func (s Store) AddSSHKey(ctx context.Context, userID, name, pub string) error {
	err := db.Q.SSHKeyAdd(ctx, db.SSHKeyAddParams{
		OwnerID:   userID,
		Name:      name,
		PublicKey: pub,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s Store) GetSSHKeys(ctx context.Context, ownerID string) ([]db.UnweaveSshKey, error) {
	keys, err := db.Q.SSHKeysGet(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	return keys, err
}
