package ssh_keys

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"

	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/tools/random"
	"golang.org/x/crypto/ssh"
)

var bitSize = 4096

type SSHKeyService struct {
	store Store
}

func NewSSHKeyService() *SSHKeyService {
	return &SSHKeyService{store: Store{}}
}

func (s *SSHKeyService) Add(ctx context.Context, userID string, params types.SSHKeyAddParams) (string, error) {
	var name string
	if params.Name != nil && *params.Name != "" {
		name = *params.Name
	}

	if name != "" {
		sshKey, err := s.store.GetSSHKeyByNameIfExists(ctx, name, userID)
		if err != nil {
			return "", fmt.Errorf("failed to get SSH key from DB: %w", err)
		}
		if sshKey != nil {
			// no-op if same key
			if sshKey.PublicKey == params.PublicKey {
				return sshKey.Name, nil
			}
			return "", &types.Error{
				Code:    http.StatusConflict,
				Message: fmt.Sprintf("Another SSH key already exists with name: %q", sshKey.Name),
			}
		}
	}

	if name == "" {
		name = "uw:" + random.GenerateRandomPhrase(4, "-") + ".pub"
	}

	err := s.store.AddSSHKey(ctx, userID, name, params.PublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to save SSH key: %w", err)
	}

	return name, nil
}

func (s *SSHKeyService) Generate(ctx context.Context, userID string, params types.SSHKeyGenerateParams) (name string, prv string, pub string, err error) {
	privateKey, publicKey, err := generateSSHKeyPair()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate SSH key pair: %w", err)
	}

	name, err = s.Add(ctx, userID, types.SSHKeyAddParams{
		Name:      params.Name,
		PublicKey: publicKey,
	})
	if err != nil {
		return "", "", "", err
	}

	return name, privateKey, publicKey, nil
}

func (s *SSHKeyService) List(ctx context.Context, userID string) ([]types.SSHKey, error) {
	sshKeys, err := s.store.GetSSHKeys(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list SSH keys from DB: %w", err)
	}

	res := make([]types.SSHKey, len(sshKeys))

	for idx, sshKey := range sshKeys {
		sshKey := sshKey
		res[idx] = types.SSHKey{
			Name:      sshKey.Name,
			PublicKey: &sshKey.PublicKey,
			CreatedAt: &sshKey.CreatedAt,
		}
	}
	return res, nil
}

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

func generatePublicKey(privateKey *rsa.PrivateKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	return ssh.MarshalAuthorizedKey(publicRsaKey), nil
}

func generateSSHKeyPair() (string, string, error) {
	privateKey, err := generatePrivateKey()
	if err != nil {
		return "", "", err
	}
	privatePEM := encodePrivateKeyToPEM(privateKey)
	publicKey, err := generatePublicKey(privateKey)
	if err != nil {
		return "", "", err
	}
	return string(privatePEM), string(publicKey), nil
}
