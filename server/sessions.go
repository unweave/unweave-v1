package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api"
	"github.com/unweave/unweave/runtime"
	"github.com/unweave/unweave/tools/random"
	"golang.org/x/crypto/ssh"
)

func setupCredentials(ctx context.Context, rt *runtime.Runtime, store *Store, sshKeyName, sshPublicKey *string) (api.SSHKey, error) {
	exists := false
	key := api.SSHKey{
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
			key.PublicKey = k.PublicKey
		}
		if err != nil && err != sql.ErrNoRows {
			return api.SSHKey{}, fmt.Errorf("failed to get ssh key from db: %w", err)
		}
	}

	if !exists && sshPublicKey != nil {
		pk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(key.PublicKey))
		if err != nil {
			return api.SSHKey{}, &api.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Invalid SSH public key",
			}
		}

		pkStr := string(ssh.MarshalAuthorizedKey(pk))
		k, err := store.SSHKey.GetByPublicKey(ctx, pkStr)
		if err == nil {
			exists = true
			key.Name = k.Name
			key.PublicKey = k.PublicKey
		}
		if err != nil && err != sql.ErrNoRows {
			return api.SSHKey{}, fmt.Errorf("failed to get ssh key from db: %w", err)
		}
	}

	if exists {
		// Key exists in the DB. Check if it exists with the provider and exit early if it
		// does.
		providerKeys, err := rt.ListSSHKeys(ctx)
		if err != nil {
			return api.SSHKey{}, fmt.Errorf("failed to list ssh keys from provider: %w", err)
		}
		for _, k := range providerKeys {
			if k.Name == key.Name {
				return key, nil
			}
		}
	}

	key, err := rt.AddSSHKey(ctx, key)
	if err != nil {
		return api.SSHKey{}, fmt.Errorf("failed to add ssh key to provider: %w", err)
	}

	if !exists {
		if err = store.SSHKey.Add(ctx, key.Name, key.PublicKey); err != nil {
			return api.SSHKey{}, fmt.Errorf("failed to add ssh key to db: %w", err)
		}
	}
	return key, nil
}

func SessionsCreate(rti runtime.Initializer, store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msgf("Executing SessionsCreate request")

		scr := api.SessionCreateRequestParams{}
		if err := render.Bind(r, &scr); err != nil {
			err = fmt.Errorf("failed to read body: %w", err)
			render.Render(w, r.WithContext(ctx), api.ErrHTTPError(err, "Invalid request body"))
			return
		}

		rt, err := rti.Initialize(ctx, scr.Provider)
		if err != nil {
			err = fmt.Errorf("failed to create runtime: %w", err)
			render.Render(w, r.WithContext(ctx), api.ErrInternalServer(err, ""))
			return
		}

		ctx = log.With().
			Stringer(ProviderCtxKey, rt.GetProvider()).
			Logger().
			WithContext(ctx)

		sshKey, err := setupCredentials(ctx, rt, store, scr.SSHKeyName, scr.SSHPublicKey)
		if err != nil {
			err = fmt.Errorf("failed to setup credentials: %w", err)
			render.Render(w, r.WithContext(ctx), api.ErrHTTPError(err, "Failed to setup credentials"))
			return
		}

		node, err := rt.InitNode(ctx, sshKey, scr.NodeTypeID, scr.Region)
		if err != nil {
			err = fmt.Errorf("failed to init node: %w", err)
			render.Render(w, r.WithContext(ctx), api.ErrHTTPError(err, "Failed to initialize node"))
			return
		}

		session := api.Session{
			SSHKey:     sshKey,
			Status:     api.StatusInitializing,
			NodeTypeID: node.TypeID,
			Region:     node.Region,
			Provider:   node.Provider,
			NodeID:     node.ID,
		}
		sessionID, err := store.Session.Add(ctx, session)
		if err != nil {
			err = fmt.Errorf("failed to create session in db: %w", err)
			render.Render(w, r.WithContext(ctx), api.ErrInternalServer(err, ""))
			return
		}
		session.ID = sessionID

		// TODO: watch status
		render.JSON(w, r, session)
	}
}

func SessionsGet(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := &api.Session{ID: uuid.New()}
		render.JSON(w, r, res)
	}
}

func SessionsList(rti runtime.Initializer, store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Ctx(ctx).Info().Msgf("Executing SessionsList request")

		sessions, err := store.Session.List(ctx)
		if err != nil {
			err = fmt.Errorf("failed to get sessions from db: %w", err)
			render.Render(w, r.WithContext(ctx), api.ErrInternalServer(err, ""))
			return
		}

		res := make([]api.Session, len(sessions))
		for idx, s := range sessions {
			s := s
			res[idx] = api.Session{
				ID: s.ID,
				SSHKey: api.SSHKey{
					Name: s.SSHKey.Name,
				},
				Status: s.Status,
			}
		}
		render.JSON(w, r, api.SessionsListResponse{Sessions: res})
	}
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
				render.Render(w, r.WithContext(ctx), &api.HTTPError{
					Code:       http.StatusNotFound,
					Message:    "Session not found",
					Suggestion: "Make sure the session id is valid",
				})
				return
			}
			err = fmt.Errorf("failed to fetch session from db %q: %w", session.ID, err)
			render.Render(w, r.WithContext(ctx), api.ErrInternalServer(err, "Failed to terminate session"))
			return
		}

		rt, err := rti.Initialize(ctx, api.RuntimeProvider(sess.Provider))
		if err != nil {
			err = fmt.Errorf("failed to create runtime %q: %w", sess.Provider, err)
			render.Render(w, r.WithContext(ctx), api.ErrInternalServer(err, "Failed to initialize runtime"))
			return
		}

		ctx = log.With().
			Stringer(ProviderCtxKey, rt.GetProvider()).
			Logger().
			WithContext(ctx)

		if err = rt.TerminateNode(ctx, sess.NodeID); err != nil {
			render.Render(w, r.WithContext(ctx), api.ErrHTTPError(err, "Failed to terminate node"))
			return
		}
		if err = store.Session.SetTerminated(ctx, session.ID); err != nil {
			log.Ctx(ctx).
				Error().
				Err(err).
				Msgf("Failed to set session %q as terminated", session.ID)
		}

		render.JSON(w, r, &api.SessionTerminateResponse{Success: true})
	}
}
