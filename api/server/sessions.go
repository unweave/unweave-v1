package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/runtime"
	"github.com/unweave/unweave/tools/random"
	"golang.org/x/crypto/ssh"
)

func setupCredentials(ctx context.Context, rt *runtime.Runtime, dbq db.Querier, userID uuid.UUID, sshKeyName, sshPublicKey *string) (types.SSHKey, error) {
	exists := false
	key := types.SSHKey{
		// This is overridden if the key already exists.
		//
		// This should most like never collide with an existing key, but it is possible.
		// In the future, we should check to see if the key already exists before
		// creating it.
		Name: "uw:" + random.GenerateRandomPhrase(4, "-"),
	}

	if sshKeyName != nil {
		k, err := dbq.SSHKeyGetByName(ctx, *sshKeyName)
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
			return types.SSHKey{}, &types.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Invalid SSH public key",
			}
		}

		pkStr := string(ssh.MarshalAuthorizedKey(pk))
		k, err := dbq.SSHKeyGetByPublicKey(ctx, pkStr)
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
		params := db.SSHKeyAddParams{
			OwnerID:   userID,
			Name:      key.Name,
			PublicKey: *key.PublicKey,
		}
		if err = dbq.SSHKeyAdd(ctx, params); err != nil {
			return types.SSHKey{}, fmt.Errorf("failed to add ssh key to db: %w", err)
		}
	}
	return key, nil
}

func SessionsCreate(rti runtime.Initializer, dbq db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := GetUserIDFromContext(ctx)
		project := GetProjectFromContext(ctx)

		zerolog.Ctx(ctx).Info().Msgf("Executing SessionsCreate request")

		scr := types.SessionCreateRequestParams{}
		if err := render.Bind(r, &scr); err != nil {
			err = fmt.Errorf("failed to read body: %w", err)
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Invalid request body"))
			return
		}

		rt, err := rti.FromAccount(ctx, userID, scr.Provider)
		if err != nil {
			err = fmt.Errorf("failed to create runtime: %w", err)
			render.Render(w, r.WithContext(ctx), ErrInternalServer(err, ""))
			return
		}

		ctx = log.With().
			Stringer(types.RuntimeProviderKey, rt.GetProvider()).
			Logger().
			WithContext(ctx)

		sshKey, err := setupCredentials(ctx, rt, dbq, userID, scr.SSHKeyName, scr.SSHPublicKey)
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

		params := db.SessionCreateParams{
			NodeID:     node.ID,
			CreatedBy:  userID,
			ProjectID:  project.ID,
			Provider:   scr.Provider.String(),
			SshKeyName: sshKey.Name,
		}
		sessionID, err := dbq.SessionCreate(ctx, params)
		if err != nil {
			err = fmt.Errorf("failed to create session in db: %w", err)
			render.Render(w, r.WithContext(ctx), ErrInternalServer(err, ""))
			return

		}

		session := &types.Session{
			ID:         sessionID,
			SSHKey:     node.KeyPair,
			Status:     types.StatusInitializing,
			NodeTypeID: node.TypeID,
			Region:     node.Region,
			Provider:   node.Provider,
		}

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

func SessionsList(rti runtime.Initializer, dbq db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		project := GetProjectFromContext(ctx)

		log.Ctx(ctx).Info().Msgf("Executing SessionsList request")

		params := db.SessionsGetParams{
			ProjectID: project.ID,
			Limit:     100,
			Offset:    0,
		}
		sessions, err := dbq.SessionsGet(ctx, params)
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
					// The generated go type for SshKeyName is a nullable string because
					// of the join, but it will never be null since session have a foreign
					// key constraint on ssh_key_id.
					Name: s.SshKeyName.String,
				},
				Status: types.DBSessionStatusToAPIStatus(s.Status),
			}
		}
		render.JSON(w, r, types.SessionsListResponse{Sessions: res})
	}
}

func SessionsTerminate(rti runtime.Initializer, dbq db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := GetUserIDFromContext(ctx)

		log.Ctx(ctx).
			Info().
			Msgf("Executing SessionsTerminate request for user %q", userID)

		session := GetSessionFromContext(ctx)
		sess, err := dbq.SessionGet(ctx, session.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r.WithContext(ctx), &types.HTTPError{
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

		rt, err := rti.FromAccount(ctx, userID, types.RuntimeProvider(sess.Provider))
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
		if err = dbq.SessionSetTerminated(ctx, session.ID); err != nil {
			log.Ctx(ctx).
				Error().
				Err(err).
				Msgf("Failed to set session %q as terminated", session.ID)
		}

		render.JSON(w, r, &types.SessionTerminateResponse{Success: true})
	}
}
