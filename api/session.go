package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	key, err := rt.AddSSHKey(ctx, key)
	if err != nil {
		return types.SSHKey{}, fmt.Errorf("failed to add ssh key to provider: %w", err)
	}

	if !exists {
		params := db.SSHKeyAddParams{
			OwnerID:   userID,
			Name:      *key.Name,
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
		userID := getUserIDFromContext(ctx)
		project := getProjectFromContext(ctx)

		logger := log.With().
			Str(ContextKeyUser, userID.String()).
			Str(ContextKeyProject, project.ID.String()).
			Logger()

		logger.Info().Msgf("Executing SessionsCreate request")

		scr := SessionCreateParams{}
		if err := render.Bind(r, &scr); err != nil {
			logger.Warn().
				Err(err).
				Stack().
				Msg("Failed to read body")

			render.Render(w, r, ErrHTTPError(err, "Invalid request body"))
			return
		}

		rt, err := rti.FromUser(userID, scr.Runtime)
		if err != nil {
			logger.Error().
				Err(err).
				Stack().
				Msg("Failed to create runtime")
			render.Render(w, r, ErrInternalServer(""))
			return
		}

		sshKey, err := setupCredentials(ctx, rt, dbq, userID, scr.SSHKey)
		if err != nil {
			logger.Error().
				Err(err).
				Stack().
				Msg("Failed to setup credentials")

			render.Render(w, r, ErrHTTPError(err, "Failed to setup credentials"))
			return
		}

		node, err := rt.InitNode(ctx, sshKey)
		if err != nil {
			logger.Warn().
				Err(err).
				Stack().
				Msg("Failed to init node")

			render.Render(w, r, ErrHTTPError(err, "Failed to initialize node"))
			return
		}

		params := db.SessionCreateParams{
			NodeID:    node.ID,
			CreatedBy: userID,
			ProjectID: project.ID,
			Runtime:   scr.Runtime.String(),
		}
		sessionID, err := dbq.SessionCreate(ctx, params)
		if err != nil {
			logger.Error().
				Err(err).
				Msg("Failed to create session")

			render.Render(w, r, ErrInternalServer(""))
			return

		}

		session := &types.Session{
			ID:     sessionID,
			SSHKey: node.KeyPair,
			Status: types.StatusInitializing,
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

type SessionsListResponse struct {
	Sessions []types.Session `json:"sessions"`
}

func SessionsList(rti runtime.Initializer, dbq db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := getUserIDFromContext(ctx)
		project := getProjectFromContext(ctx)

		logger := log.With().
			Str(ContextKeyUser, userID.String()).
			Str(ContextKeyProject, project.ID.String()).
			Logger()

		logger.Info().Msgf("Executing SessionsList request")
		params := db.SessionsGetParams{
			ProjectID: project.ID,
			Limit:     100,
			Offset:    0,
		}
		sessions, err := dbq.SessionsGet(ctx, params)
		if err != nil {
			logger.Error().
				Err(err).
				Msg("Failed to get sessions")

			render.Render(w, r, ErrInternalServer(""))
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
					Name: &s.SshKeyName.String,
				},
				Status: dbSessionStatusToAPIStatus(s.Status),
			}
		}
		render.JSON(w, r, SessionsListResponse{Sessions: res})
	}
}

type SessionTerminateResponse struct {
	Success bool `json:"success"`
}

func SessionsTerminate(rti runtime.Initializer, dbq db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := getUserIDFromContext(ctx)
		logger := log.With().Str(ContextKeyUser, userID.String()).Logger()

		logger.Info().Msgf("Executing SessionsTerminate request for user %q", userID)

		// fetch from url params and try converting to uuid
		id := chi.URLParam(r, "sessionID")
		sessionID, err := uuid.Parse(id)
		if err != nil {
			render.Render(w, r, &HTTPError{
				Code:       400,
				Message:    "Invalid session id",
				Suggestion: "Make sure the session id is a valid UUID",
			})
			return
		}

		sess, err := dbq.SessionGet(ctx, sessionID)
		if err != nil {
			if err == sql.ErrNoRows {
				render.Render(w, r, &HTTPError{
					Code:       404,
					Message:    "Session not found",
					Suggestion: "Make sure the session id is valid",
				})
				return
			}
			logger.Error().
				Err(err).
				Msgf("Error fetching session %q", sessionID)

			render.Render(w, r, ErrInternalServer("Failed to terminate session"))
			return
		}

		rt, err := rti.FromUser(sessionID, types.RuntimeProvider(sess.Runtime))
		if err != nil {
			logger.Error().
				Err(err).
				Msg("Failed to create runtime" + sess.Runtime)

			render.Render(w, r, ErrInternalServer("Failed to initialize runtime"))
			return
		}

		if err = rt.TerminateNode(ctx, sess.NodeID); err != nil {
			render.Render(w, r, ErrHTTPError(err, "Failed to terminate node"))
			return
		}
		if err = dbq.SessionSetTerminated(ctx, sessionID); err != nil {
			logger.Error().
				Err(err).
				Msgf("Failed to set session %q as terminated", sessionID)
		}

		render.JSON(w, r, &SessionTerminateResponse{Success: true})
	}
}
