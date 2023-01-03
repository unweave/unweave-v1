package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/runtime"
	"github.com/unweave/unweave/types"
	"golang.org/x/crypto/ssh"
)

type SessionCreateParams struct {
	Runtime types.RuntimeProvider `json:"runtime"`
	TypeID  *string               `json:"typeID,omitempty"`
	Region  *string               `json:"region,omitempty"`
	SSHKey  *types.SSHKey         `json:"sshKey"`
	Specs   *types.NodeSpecs      `json:"specs,omitempty"`
}

func (s *SessionCreateParams) Bind(r *http.Request) error {
	if s.Runtime == "" {
		return &HTTPError{
			Code:       400,
			Message:    "Invalid request body: field 'runtime' is required",
			Suggestion: fmt.Sprintf("Use %q or %q as the runtime provider", types.LambdaLabsProvider, types.UnweaveProvider),
		}
	}
	if s.Runtime != types.LambdaLabsProvider && s.Runtime != types.UnweaveProvider {
		return &HTTPError{
			Code:       400,
			Message:    "Invalid runtime provider: " + string(s.Runtime),
			Suggestion: fmt.Sprintf("Use %q or %q as the runtime provider", types.LambdaLabsProvider, types.UnweaveProvider),
		}
	}
	return nil
}

func setupCredentials(ctx context.Context, rt *runtime.Runtime, dbq db.Querier, userID uuid.UUID, sshKey *types.SSHKey) (types.SSHKey, error) {
	exists := false

	key := types.SSHKey{}
	if sshKey != nil {
		key = *sshKey
	}

	if key.Name != nil {
		k, err := dbq.SSHKeyGetByName(ctx, *key.Name)
		if err == nil {
			exists = true
			key.PublicKey = &k.PublicKey
		}
		if err != nil && err != sql.ErrNoRows {
			return types.SSHKey{}, fmt.Errorf("failed to get ssh key from db: %w", err)
		}
	}

	if !exists && key.PublicKey != nil {
		pk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(*key.PublicKey))
		if err != nil {
			return types.SSHKey{}, &HTTPError{
				Code:    400,
				Message: "Invalid SSH public key",
			}
		}

		pkStr := string(ssh.MarshalAuthorizedKey(pk))
		k, err := dbq.SSHKeyGetByPublicKey(ctx, pkStr)
		if err == nil {
			exists = true
			key.Name = &k.Name
		}
		if err != nil && err != sql.ErrNoRows {
			return types.SSHKey{}, fmt.Errorf("failed to get ssh key from db: %w", err)
		}
	}

	if exists {
		providerKeys, err := rt.ListSSHKeys(ctx)
		if err != nil {
			return types.SSHKey{}, fmt.Errorf("failed to list ssh keys from provider: %w", err)
		}
		for _, k := range providerKeys {
			if *k.Name == *key.Name {
				return key, nil
			}
		}
	}

	if _, err := rt.AddSSHKey(ctx, key); err != nil {
		return types.SSHKey{}, fmt.Errorf("failed to add ssh key to provider: %w", err)
	}

	if !exists {
		params := db.SSHKeyAddParams{
			OwnerID:   userID,
			Name:      *key.Name,
			PublicKey: *key.PublicKey,
		}
		if err := dbq.SSHKeyAdd(ctx, params); err != nil {
			return types.SSHKey{}, fmt.Errorf("failed to add ssh key to db: %w", err)
		}
	}
	return key, nil
}

func SessionsCreate(rti runtime.Initializer, dbq db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		scr := SessionCreateParams{}
		if err := render.Bind(r, &scr); err != nil {
			log.Warn().
				Err(err).
				Msg("failed to read body")

			render.Render(w, r, ErrHTTPError(err, "Invalid request body"))
			return
		}

		rt, err := rti.FromUser(uuid.New(), scr.Runtime)
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to create runtime")
			render.Render(w, r, ErrInternalServer(""))
		}

		sshKey, err := setupCredentials(ctx, rt, dbq, uuid.New(), scr.SSHKey)
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to setup credentials")

			render.Render(w, r, ErrHTTPError(err, "Failed to setup credentials"))
			return
		}

		node, err := rt.InitNode(ctx, sshKey)
		if err != nil {
			log.Warn().
				Err(err).
				Msg("failed to init node")

			render.Render(w, r, ErrHTTPError(err, "Failed to initialize node"))
			return
		}

		// add to db
		res := &types.Session{ID: node.ID, SSHKey: node.KeyPair}

		// watch status
		render.JSON(w, r, res)
	}
}
