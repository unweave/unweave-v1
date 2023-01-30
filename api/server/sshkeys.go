package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/rs/zerolog/log"
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

func saveSSHKey(ctx context.Context, dbq db.Querier, userID uuid.UUID, name, publicKey string) error {
	arg := db.SSHKeyAddParams{
		OwnerID:   userID,
		Name:      name,
		PublicKey: publicKey,
	}
	err := dbq.SSHKeyAdd(ctx, arg)
	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) {
			// We already check the unique constraint on the name column, so this
			// should only happen if the public key is a duplicate.
			if e.Code == pgerrcode.UniqueViolation {
				return &types.HTTPError{
					Code:    http.StatusConflict,
					Message: "Public key already exists",
					Suggestion: "Public keys in Unweave have to be globally unique. " +
						"It could be that you added this key earlier or that you " +
						"added it to another account. If you've already added this " +
						"key to your account, remove it first.",
				}
			}
		}
		return fmt.Errorf("failed to add ssh key to db: %w", err)
	}
	return nil
}

// SSHKeyAdd adds an SSH key to the user's account.
//
// This does not add the key to the user's configured providers. That is done lazily
// when the user first tries to use the key.
func SSHKeyAdd(dbq db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := GetUserIDFromContext(ctx)

		log.Ctx(ctx).Info().Msgf("Executing SSHKeyAdd request")

		params := types.SSHKeyAddParams{}
		if err := render.Bind(r, &params); err != nil {
			err = fmt.Errorf("failed to read body: %w", err)
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Invalid request body"))
			return
		}

		if params.Name != nil {
			p := db.SSHKeyGetByNameParams{Name: *params.Name, OwnerID: userID}
			k, err := dbq.SSHKeyGetByName(ctx, p)
			if err == nil {
				render.Render(w, r.WithContext(ctx), &types.HTTPError{
					Code:    http.StatusNotFound,
					Message: fmt.Sprintf("SSH key already exists with name: %q", k.Name),
				})
				return
			}
			if err != nil && err != sql.ErrNoRows {
				err = fmt.Errorf("failed to get ssh key from db: %w", err)
				render.Render(w, r.WithContext(ctx), ErrInternalServer(err, ""))
				return
			}
		} else {
			// This should most like never collide with an existing key, but it is possible.
			// In the future, we should check to see if the key already exists before
			// creating it.
			params.Name = tools.Stringy("uw:" + random.GenerateRandomPhrase(4, "-"))
		}

		if err := saveSSHKey(ctx, dbq, userID, *params.Name, params.PublicKey); err != nil {
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to save SSH key"))
			return
		}
		render.JSON(w, r, &types.SSHKeyAddResponse{Success: true})
	}
}

func SSHKeyGenerate(dbq db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := GetUserIDFromContext(ctx)

		log.Ctx(ctx).Info().Msgf("Executing SSHKeyCreate request")

		privateKey, publicKey, err := createSSHKeyPair()
		if err != nil {
			err = fmt.Errorf("failed to generate ssh key pair: %w", err)
			render.Render(w, r.WithContext(ctx), ErrInternalServer(err, ""))
			return
		}

		params := types.SSHKeyGenerateParams{}
		render.Bind(r, &params)

		name := "uw:" + random.GenerateRandomPhrase(4, "-")
		if params.Name != nil {
			name = *params.Name
		}

		arg := db.SSHKeyAddParams{
			OwnerID:   userID,
			Name:      name,
			PublicKey: publicKey,
		}
		if err = dbq.SSHKeyAdd(ctx, arg); err != nil {
			err = fmt.Errorf("failed to add ssh key to db: %w", err)
			render.Render(w, r.WithContext(ctx), ErrInternalServer(err, ""))
			return
		}

		res := types.SSHKeyGenerateResponse{
			Name:       arg.Name,
			PublicKey:  publicKey,
			PrivateKey: privateKey,
		}
		render.JSON(w, r, &res)
	}
}

func SSHKeyList(dbq db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := GetUserIDFromContext(ctx)

		log.Ctx(ctx).Info().Msgf("Executing SSHKeyList request")

		keys, err := dbq.SSHKeysGet(ctx, userID)
		if err != nil {
			err = fmt.Errorf("failed to list ssh keys from db: %w", err)
			render.Render(w, r.WithContext(ctx), ErrInternalServer(err, ""))
			return
		}

		res := types.SSHKeyListResponse{
			Keys: make([]types.SSHKey, len(keys)),
		}
		for idx, key := range keys {
			res.Keys[idx] = types.SSHKey{
				Name:      key.Name,
				PublicKey: &key.PublicKey,
				CreatedAt: &key.CreatedAt,
			}
		}
		render.JSON(w, r, res)
	}
}
