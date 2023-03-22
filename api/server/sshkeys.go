package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"net/http"

	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/tools"
	"github.com/unweave/unweave/tools/random"
	"golang.org/x/crypto/ssh"
)

var bitSize = 4096

func generatePrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	der := x509.MarshalPKCS1PrivateKey(privateKey)
	block := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   der,
	}
	privatePEM := pem.EncodeToMemory(&block)
	return privatePEM
}

func generatePublicKey(privateKey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privateKey)
	if err != nil {
		return nil, err
	}
	return ssh.MarshalAuthorizedKey(publicRsaKey), nil
}

func createSSHKeyPair() (string, string, error) {
	privateKey, err := generatePrivateKey()
	if err != nil {
		return "", "", err
	}
	privatePEM := encodePrivateKeyToPEM(privateKey)
	publicKey, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	return string(privatePEM), string(publicKey), nil
}

func saveSSHKey(ctx context.Context, accountID string, name, publicKey string) error {
	arg := db.SSHKeyAddParams{
		OwnerID:   accountID,
		Name:      name,
		PublicKey: publicKey,
	}
	err := db.Q.SSHKeyAdd(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to add ssh key to db: %w", err)
	}
	return nil
}

type SSHKeyService struct {
	srv *Service
}

func (s *SSHKeyService) Add(ctx context.Context, params types.SSHKeyAddParams) error {
	if params.Name != nil {
		p := db.SSHKeyGetByNameParams{Name: *params.Name, OwnerID: s.srv.cid}
		k, err := db.Q.SSHKeyGetByName(ctx, p)
		if err == nil {
			// Check if it is the same key. If yes, then this is a no-op.
			if k.PublicKey == params.PublicKey {
				return nil
			}
			return &types.Error{
				Code:    http.StatusConflict,
				Message: fmt.Sprintf("SSH key already exists with name: %q", k.Name),
			}
		}
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("failed to get ssh key from db: %w", err)
		}
	} else {
		// This should most like never collide with an existing key, but it is possible.
		// In the future, we should check to see if the key already exists before
		// creating it.
		params.Name = tools.Stringy("uw:" + random.GenerateRandomPhrase(4, "-"))
	}

	if err := saveSSHKey(ctx, s.srv.cid, *params.Name, params.PublicKey); err != nil {
		return fmt.Errorf("failed to save ssh key: %w", err)
	}
	return nil
}

func (s *SSHKeyService) Generate(ctx context.Context, params types.SSHKeyGenerateParams) (name string, prv string, pub string, err error) {
	privateKey, publicKey, err := createSSHKeyPair()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate ssh key pair: %w", err)
	}

	name = "uw:" + random.GenerateRandomPhrase(4, "-")
	if params.Name != nil {
		name = *params.Name
	}

	arg := db.SSHKeyAddParams{
		OwnerID:   s.srv.cid,
		Name:      name,
		PublicKey: publicKey,
	}
	if err = db.Q.SSHKeyAdd(ctx, arg); err != nil {
		return "", "", "", fmt.Errorf("failed to add ssh key to db: %w", err)
	}
	return name, privateKey, publicKey, nil
}

func (s *SSHKeyService) List(ctx context.Context) ([]types.SSHKey, error) {
	keys, err := db.Q.SSHKeysGet(ctx, s.srv.cid)
	if err != nil {
		return nil, fmt.Errorf("failed to list ssh keys from db: %w", err)
	}

	res := make([]types.SSHKey, len(keys))

	for idx, key := range keys {
		key := key
		res[idx] = types.SSHKey{
			Name:      key.Name,
			PublicKey: &key.PublicKey,
			CreatedAt: &key.CreatedAt,
		}
	}
	return res, nil
}
