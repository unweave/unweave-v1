package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/runtime"
	"github.com/unweave/unweave/tools/random"
	"github.com/unweave/unweave/types"
	"golang.org/x/crypto/ssh"
)

type SessionCreateParams struct {
	Provider     types.RuntimeProvider `json:"provider"`
	NodeTypeID   string                `json:"nodeTypeID,omitempty"`
	Region       *string               `json:"region,omitempty"`
	SSHKeyName   *string               `json:"sshKeyName"`
	SSHPublicKey *string               `json:"sshPublicKey"`
}

func (s *SessionCreateParams) Bind(r *http.Request) error {
	if s.Provider == "" {
		return &HTTPError{
			Code:       http.StatusBadRequest,
			Message:    "Invalid request body: field 'runtime' is required",
			Suggestion: fmt.Sprintf("Use %q or %q as the runtime provider", types.LambdaLabsProvider, types.UnweaveProvider),
		}
	}
	if s.Provider != types.LambdaLabsProvider && s.Provider != types.UnweaveProvider {
		return &HTTPError{
			Code:       http.StatusBadRequest,
			Message:    "Invalid runtime provider: " + string(s.Provider),
			Suggestion: fmt.Sprintf("Use %q or %q as the runtime provider", types.LambdaLabsProvider, types.UnweaveProvider),
		}
	}
	return nil
}

func setupCredentials(ctx context.Context, rt *runtime.Runtime, store *Store, sshKeyName, sshPublicKey *string) (types.SSHKey, error) {
	exists := false
	key := types.SSHKey{
		// This is overridden if the key already exists.
		//
		// This should most like never collide with an existing key, but it is possible.
		// In the future, we should check to see if the key already exists before creating
		// it.
		Name: "uw:" + random.GenerateRandomPhrase(4, "-"),
	}

	if sshKeyName != nil {
		k, err := store.SSHKey.GetByName(ctx, *sshKeyName)
		if err == nil {
			exists = true
			key.Name = *sshKeyName
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
				Code:    http.StatusBadRequest,
				Message: "Invalid SSH public key",
			}
		}

		pkStr := string(ssh.MarshalAuthorizedKey(pk))
		k, err := store.SSHKey.GetByPublicKey(ctx, pkStr)
		if err == nil {
			exists = true
			key.Name = k.Name
			key.PublicKey = &k.PublicKey
		}
		if err != nil && err != sql.ErrNoRows {
			return types.SSHKey{}, fmt.Errorf("failed to get ssh key from db: %w", err)
		}
	}

	if exists {
		// Key exists in the DB. Check if it exists with the provider and exit early if it
		// does.
		providerKeys, err := rt.ListSSHKeys(ctx)
		if err != nil {
			return types.SSHKey{}, fmt.Errorf("failed to list ssh keys from provider: %w", err)
		}
		for _, k := range providerKeys {
			if k.Name == key.Name {
				return key, nil
			}
		}
	}

	key, err := rt.AddSSHKey(ctx, key)
	if err != nil {
		return types.SSHKey{}, fmt.Errorf("failed to add ssh key to provider: %w", err)
	}

	if !exists {
		if err = store.SSHKey.Add(ctx, key.Name, *key.PublicKey); err != nil {
			return types.SSHKey{}, fmt.Errorf("failed to add ssh key to db: %w", err)
		}
	}
	return key, nil
}

type SessionCreateResponse struct {
	Session types.Session `json:"session"`
}

func SessionsCreate(rti runtime.Initializer, store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msgf("Executing SessionsCreate request")

		scr := SessionCreateParams{}
		if err := render.Bind(r, &scr); err != nil {
			err = fmt.Errorf("failed to read body: %w", err)
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Invalid request body"))
			return
		}

		rt, err := rti.Initialize(ctx, scr.Provider)
		if err != nil {
			err = fmt.Errorf("failed to create runtime: %w", err)
			render.Render(w, r.WithContext(ctx), ErrInternalServer(err, ""))
			return
		}

		ctx = log.With().
			Stringer(types.RuntimeProviderKey, rt.GetProvider()).
			Logger().
			WithContext(ctx)

		sshKey, err := setupCredentials(ctx, rt, store, scr.SSHKeyName, scr.SSHPublicKey)
		if err != nil {
			err = fmt.Errorf("failed to setup credentials: %w", err)
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to setup credentials"))
			return
		}

		node, err := rt.InitNode(ctx, sshKey, scr.NodeTypeID, scr.Region)
		if err != nil {
			err = fmt.Errorf("failed to init node: %w", err)
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to initialize node"))
			return
		}

		session := types.Session{
			SSHKey:     node.KeyPair,
			Status:     types.StatusInitializing,
			NodeTypeID: node.TypeID,
			Region:     node.Region,
			Provider:   node.Provider,
		}
		sessionID, err := store.Session.Add(ctx, session)
		if err != nil {
			err = fmt.Errorf("failed to create session in db: %w", err)
			render.Render(w, r.WithContext(ctx), ErrInternalServer(err, ""))
			return
		}
		session.ID = sessionID

		// TODO: watch status
		render.JSON(w, r, session)
	}
}

func SessionsGet(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := &types.Session{ID: uuid.New()}
		render.JSON(w, r, res)
	}
}

type SessionsListResponse struct {
	Sessions []types.Session `json:"sessions"`
}

func SessionsList(rti runtime.Initializer, store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msgf("Executing SessionsList request")

		sessions, err := store.Session.List(ctx)
		if err != nil {
			err = fmt.Errorf("failed to get sessions from db: %w", err)
			render.Render(w, r.WithContext(ctx), ErrInternalServer(err, ""))
			return
		}

		res := make([]types.Session, len(sessions))
		for idx, s := range sessions {
			s := s
			res[idx] = types.Session{
				ID: s.ID,
				SSHKey: types.SSHKey{
					Name: s.SSHKey.Name,
				},
				Status: s.Status,
			}
		}
		render.JSON(w, r, SessionsListResponse{Sessions: res})
	}
}

type SessionTerminateResponse struct {
	Success bool `json:"success"`
}

func SessionsTerminate(rti runtime.Initializer, store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		log.Ctx(ctx).
			Info().
			Msgf("Executing SessionsTerminate request")

		session := GetSessionFromContext(ctx)
		sess, err := store.Session.Get(ctx, session.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r.WithContext(ctx), &HTTPError{
					Code:       http.StatusNotFound,
					Message:    "Session not found",
					Suggestion: "Make sure the session id is valid",
				})
				return
			}
			err = fmt.Errorf("failed to fetch session from db %q: %w", session.ID, err)
			render.Render(w, r.WithContext(ctx), ErrInternalServer(err, "Failed to terminate session"))
			return
		}

		rt, err := rti.Initialize(ctx, types.RuntimeProvider(sess.Provider))
		if err != nil {
			err = fmt.Errorf("failed to create runtime %q: %w", sess.Provider, err)
			render.Render(w, r.WithContext(ctx), ErrInternalServer(err, "Failed to initialize runtime"))
			return
		}

		ctx = log.With().
			Stringer(types.RuntimeProviderKey, rt.GetProvider()).
			Logger().
			WithContext(ctx)

		if err = rt.TerminateNode(ctx, sess.NodeID); err != nil {
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to terminate node"))
			return
		}
		if err = store.Session.SetTerminated(ctx, session.ID); err != nil {
			log.Ctx(ctx).
				Error().
				Err(err).
				Msgf("Failed to set session %q as terminated", session.ID)
		}

		render.JSON(w, r, &SessionTerminateResponse{Success: true})
	}
}
