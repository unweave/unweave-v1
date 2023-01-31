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
	"github.com/unweave/unweave/tools"
	"github.com/unweave/unweave/tools/random"
	"golang.org/x/crypto/ssh"
)

func registerCredentials(ctx context.Context, rt *runtime.Runtime, key types.SSHKey) error {
	// Check if it exists with the provider and exit early if it does
	providerKeys, err := rt.ListSSHKeys(ctx)
	if err != nil {
		return fmt.Errorf("failed to list ssh keys from provider: %w", err)
	}
	for _, k := range providerKeys {
		if k.Name == key.Name {
			return nil
		}
	}
	if _, err = rt.AddSSHKey(ctx, key); err != nil {
		return fmt.Errorf("failed to add ssh key to provider: %w", err)
	}
	return nil
}

func fetchCredentials(ctx context.Context, userID uuid.UUID, sshKeyName, sshPublicKey *string) (types.SSHKey, error) {
	if sshKeyName == nil && sshPublicKey == nil {
		return types.SSHKey{}, &types.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Either Key name or Public Key must be provided",
		}
	}

	if sshKeyName != nil {
		params := db.SSHKeyGetByNameParams{Name: *sshKeyName, OwnerID: userID}
		k, err := db.Q.SSHKeyGetByName(ctx, params)
		if err == nil {
			return types.SSHKey{
				Name:      k.Name,
				PublicKey: &k.PublicKey,
				CreatedAt: &k.CreatedAt,
			}, nil
		}
		if err != sql.ErrNoRows {
			return types.SSHKey{}, &types.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get SSH key",
				Err:     fmt.Errorf("failed to get ssh key from db: %w", err),
			}
		}
	}

	// Not found by name, try public key
	if sshPublicKey != nil {
		pk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(*sshPublicKey))
		if err != nil {
			return types.SSHKey{}, &types.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Invalid SSH public key",
			}
		}

		pkStr := string(ssh.MarshalAuthorizedKey(pk))
		params := db.SSHKeyGetByPublicKeyParams{PublicKey: pkStr, OwnerID: userID}
		k, err := db.Q.SSHKeyGetByPublicKey(ctx, params)
		if err == nil {
			return types.SSHKey{
				Name:      k.Name,
				PublicKey: &k.PublicKey,
				CreatedAt: &k.CreatedAt,
			}, nil
		}
		if err != sql.ErrNoRows {
			return types.SSHKey{}, &types.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get SSH key",
				Err:     fmt.Errorf("failed to get ssh key from db: %w", err),
			}
		}
	}

	if sshPublicKey == nil {
		return types.SSHKey{}, &types.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "SSH key not found",
		}
	}
	if sshKeyName == nil {
		sshKeyName = tools.Stringy("uw:" + random.GenerateRandomPhrase(4, "-"))
	}

	// Key doesn't exist in db, but the user provided a public key, so add it to the db
	if err := saveSSHKey(ctx, userID, *sshKeyName, *sshPublicKey); err != nil {
		return types.SSHKey{}, &types.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to save SSH key",
		}
	}
	return types.SSHKey{
		Name:      *sshKeyName,
		PublicKey: sshPublicKey,
	}, nil
}

func SessionsCreate(rti runtime.Initializer) http.HandlerFunc {
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

		rt, err := rti.Initialize(ctx, userID, scr.Provider, scr.ProviderToken)
		if err != nil {
			err = fmt.Errorf("failed to create runtime: %w", err)
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to initialize provider"))
			return
		}

		ctx = log.With().
			Stringer(types.RuntimeProviderKey, rt.GetProvider()).
			Logger().
			WithContext(ctx)

		sshKey, err := fetchCredentials(ctx, userID, scr.SSHKeyName, scr.SSHPublicKey)
		if err != nil {
			err = fmt.Errorf("failed to setup credentials: %w", err)
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to setup credentials"))
			return
		}
		if err = registerCredentials(ctx, rt, sshKey); err != nil {
			err = fmt.Errorf("failed to register credentials: %w", err)
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to register credentials"))
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
		sessionID, err := db.Q.SessionCreate(ctx, params)
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

func SessionsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	project := GetProjectFromContext(ctx)

	log.Ctx(ctx).Info().Msgf("Executing SessionsList request")

	params := db.SessionsGetParams{
		ProjectID: project.ID,
		Limit:     100,
		Offset:    0,
	}
	sessions, err := db.Q.SessionsGet(ctx, params)
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

func SessionsTerminate(rti runtime.Initializer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := GetUserIDFromContext(ctx)

		log.Ctx(ctx).
			Info().
			Msgf("Executing SessionsTerminate request for user %q", userID)

		session := GetSessionFromContext(ctx)
		sess, err := db.Q.SessionGet(ctx, session.ID)
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

		provider := types.RuntimeProvider(sess.Provider)
		str := types.SessionTerminateRequestParams{}
		render.Bind(r, &str)

		var rt *runtime.Runtime

		rt, err = rti.Initialize(ctx, userID, provider, str.ProviderToken)
		if err != nil {
			err = fmt.Errorf("failed to create runtime %q: %w", provider, err)
			render.Render(w, r.WithContext(ctx), ErrHTTPError(err, "Failed to initialize runtime"))
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
		if err = db.Q.SessionSetTerminated(ctx, session.ID); err != nil {
			log.Ctx(ctx).
				Error().
				Err(err).
				Msgf("Failed to set session %q as terminated", session.ID)
		}

		render.JSON(w, r, &types.SessionTerminateResponse{Success: true})
	}
}
