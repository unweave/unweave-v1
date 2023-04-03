package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/runtime"
	"github.com/unweave/unweave/tools"
	"github.com/unweave/unweave/tools/random"
	"golang.org/x/crypto/ssh"
)

type ConnectionInfoV1 struct {
	Version int    `json:"version"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	User    string `json:"user"`
}

func handleSessionError(ctx context.Context, sessionID string, err error, msg string) {
	ctx, _ = context.WithCancel(ctx) // make sure this doesn't fail because of a parent cancelled context

	var e *types.Error
	if errors.As(err, &e) {
		msg += ": " + e.Message
	}
	msg += ": " + err.Error()

	log.Ctx(ctx).Error().Err(err).Msg(msg)

	params := db.SessionSetErrorParams{
		ID: sessionID,
		Error: sql.NullString{
			String: msg,
			Valid:  true,
		},
	}
	if err = db.Q.SessionSetError(ctx, params); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to set session error")
	}
}

func registerCredentials(ctx context.Context, rt *runtime.Runtime, keys []types.SSHKey) error {
	// Check if it exists with the provider and exit early if it does
	providerKeys, err := rt.Node.ListSSHKeys(ctx)
	if err != nil {
		return fmt.Errorf("failed to list ssh keys from provider: %w", err)
	}

	for _, key := range keys {
		for _, k := range providerKeys {
			if k.Name == key.Name {
				return nil
			}
		}
		if _, err = rt.Node.AddSSHKey(ctx, key); err != nil {
			return fmt.Errorf("failed to add ssh key to provider: %w", err)
		}
	}

	return nil
}

func fetchCredentials(ctx context.Context, userID string, sshKeyName, sshPublicKey *string) (types.SSHKey, error) {
	if sshKeyName == nil && sshPublicKey == nil {
		return types.SSHKey{}, &types.Error{
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
			return types.SSHKey{}, &types.Error{
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
			return types.SSHKey{}, &types.Error{
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
			return types.SSHKey{}, &types.Error{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get SSH key",
				Err:     fmt.Errorf("failed to get ssh key from db: %w", err),
			}
		}
	}

	// Public key wasn't provided	 and key wasn't found by name
	if sshPublicKey == nil {
		return types.SSHKey{}, &types.Error{
			Code:    http.StatusBadRequest,
			Message: "SSH key not found",
		}
	}
	if sshKeyName == nil || *sshKeyName == "" {
		sshKeyName = tools.Stringy("uw:" + random.GenerateRandomPhrase(4, "-"))
	}

	// Key doesn't exist in db, but the user provided a public key, so add it to the db
	if err := saveSSHKey(ctx, userID, *sshKeyName, *sshPublicKey); err != nil {
		return types.SSHKey{}, &types.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to save SSH key",
		}
	}
	return types.SSHKey{
		Name:      *sshKeyName,
		PublicKey: sshPublicKey,
	}, nil
}

func updateConnectionInfo(ctx context.Context, rt runtime.Node, nodeID string, sessionID string) error {
	ctx, _ = context.WithCancel(ctx) // make sure this doesn't fail because of a parent cancelled context
	connInfo, err := rt.GetConnectionInfo(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("failed to get connection info: %w", err)
	}

	connInfoJSON, err := json.Marshal(connInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal connection info: %w", err)
	}
	params := db.SessionUpdateConnectionInfoParams{
		ID:             sessionID,
		ConnectionInfo: connInfoJSON,
	}
	if e := db.Q.SessionUpdateConnectionInfo(ctx, params); e != nil {
		return fmt.Errorf("failed to update connection info: %w", e)
	}
	return nil
}

func updateExecStatus(ctx context.Context, execID string, status types.NodeStatus) {
	ctx, _ = context.WithCancel(ctx) // make sure this doesn't fail because of a parent cancelled context
	params := db.SessionStatusUpdateParams{
		ID:     execID,
		Status: db.UnweaveSessionStatus(status),
	}
	if err := db.Q.SessionStatusUpdate(ctx, params); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to update exec status")
	}
}

type ExecService struct {
	srv *Service
}

func (s *ExecService) Create(ctx context.Context, projectID string, params types.ExecCreateParams) (*types.Exec, error) {
	rt, err := s.srv.InitializeRuntime(ctx, params.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create runtime: %w", err)
	}

	ctx = log.With().
		Stringer(types.RuntimeProviderKey, rt.Node.GetProvider()).
		Logger().
		WithContext(ctx)

	userKey, err := fetchCredentials(ctx, s.srv.cid, params.SSHKeyName, params.SSHPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to setup credentials: %w", err)
	}
	prv, pub, err := generateSSHKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ssh key pair: %w", err)
	}

	adminKey := types.SSHKey{
		Name:      "umk-" + random.GenerateRandomAdjectiveNounTriplet(),
		PublicKey: &pub,
		CreatedAt: nil,
	}
	keys := []types.SSHKey{userKey, adminKey}

	if err = registerCredentials(ctx, rt, keys); err != nil {
		return nil, fmt.Errorf("failed to register credentials: %w", err)
	}

	node, err := rt.Node.InitNode(ctx, keys, params.NodeTypeID, params.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to init node: %w", err)
	}
	if _, err := s.srv.vault.SetSecret(ctx, prv, &node.ID); err != nil {
		return nil, fmt.Errorf("failed to store private key: %w", err)
	}

	specs, err := json.Marshal(node.Specs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal connection info: %w", err)
	}

	sshKey, err := db.Q.SSHKeysGet(ctx, s.srv.cid)
	if err != nil {
		return nil, fmt.Errorf("failed to get ssh key ids: %w", err)
	}

	np := db.NodeCreateParams{
		ID:        node.ID,
		Provider:  string(rt.Node.GetProvider()),
		Region:    node.Region,
		Spec:      specs,
		Status:    string(types.StatusInitializing),
		OwnerID:   s.srv.aid,
		SshKeyIds: []string{sshKey[0].ID},
	}
	if err = db.Q.NodeCreate(ctx, np); err != nil {
		return nil, fmt.Errorf("failed to create node in db: %w", err)
	}

	connInfo, err := json.Marshal(ConnectionInfoV1{Version: 1})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal connection info: %w", err)
	}

	// Set commit details if provided
	var command []string
	var commitID, gitRemoteURL sql.NullString

	if params.Ctx.Command != nil {
		command = params.Ctx.Command
	}
	if params.Ctx.CommitID != nil {
		commitID = sql.NullString{String: *params.Ctx.CommitID, Valid: true}
	}
	if params.Ctx.GitURL != nil {
		gitRemoteURL = sql.NullString{String: *params.Ctx.GitURL, Valid: true}
	}

	dbp := db.SessionCreateParams{
		NodeID:         node.ID,
		CreatedBy:      s.srv.cid,
		ProjectID:      projectID,
		Region:         node.Region,
		Name:           random.GenerateRandomPhrase(4, "-"),
		ConnectionInfo: connInfo,
		CommitID:       commitID,
		GitRemoteUrl:   gitRemoteURL,
		Command:        command,
		SshKeyName:     userKey.Name,
	}
	execID, err := db.Q.SessionCreate(ctx, dbp)
	if err != nil {
		return nil, fmt.Errorf("failed to create session in db: %w", err)
	}
	ctx = log.With().Str(SessionIDCtxKey, execID).Logger().WithContext(ctx)

	imageURI := "alpine:latest"
	if params.Ctx.BuildID != nil {
		var err error
		imageURI, err = s.srv.Builder.GetImageURI(ctx, *params.Ctx.BuildID)
		if err != nil {
			log.Error().Err(err).Msgf("failed to get image uri: %w", err)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get image uri: %w", err)
	}
	if err := rt.Session.Init(ctx, node, []types.SSHKey{userKey}, imageURI); err != nil {
		go handleSessionError(ctx, execID, err, "failed to init session")
		return nil, fmt.Errorf("failed to init session: %w", err)
	}
	if err := rt.Session.Exec(ctx, node.ID, execID, params.Ctx, true); err != nil {
		go handleSessionError(ctx, execID, err, "failed to run exec")
		return nil, fmt.Errorf("failed to to run exec: %w", err)
	}

	createdAt := time.Now()
	session := &types.Exec{
		ID:         execID,
		Name:       dbp.Name,
		SSHKey:     node.KeyPair,
		Connection: nil,
		Status:     types.StatusInitializing,
		CreatedAt:  &createdAt,
		NodeTypeID: node.TypeID,
		Region:     node.Region,
		Provider:   node.Provider,
		Ctx: types.ExecCtx{
			Command:  command,
			CommitID: &commitID.String,
			GitURL:   &gitRemoteURL.String,
			BuildID:  params.Ctx.BuildID,
		},
	}

	return session, nil
}

func (s *ExecService) Get(ctx context.Context, sessionID string) (*types.Exec, error) {
	dbs, err := db.Q.MxSessionGet(ctx, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &types.Error{
				Code:    http.StatusNotFound,
				Message: "Session not found",
			}
		}
		return nil, fmt.Errorf("failed to get session from db: %w", err)
	}

	connInfo := &ConnectionInfoV1{}
	if err := json.Unmarshal(dbs.ConnectionInfo, connInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal connection info: %w", err)
	}

	session := &types.Exec{
		ID:   sessionID,
		Name: dbs.Name,
		SSHKey: types.SSHKey{
			Name:      dbs.SshKeyName,
			PublicKey: &dbs.PublicKey,
			CreatedAt: &dbs.SshKeyCreatedAt,
		},
		Connection: &types.ConnectionInfo{
			Host: connInfo.Host,
			Port: connInfo.Port,
			User: connInfo.User,
		},
		Status:     types.NodeStatus(dbs.Status),
		CreatedAt:  &dbs.CreatedAt,
		NodeTypeID: dbs.NodeID,
		Region:     dbs.Region,
		Provider:   types.Provider(dbs.Provider),
		Ctx:        types.ExecCtx{},
	}
	return session, nil
}

func (s *ExecService) List(ctx context.Context, projectID string, listTerminated bool) ([]types.Exec, error) {
	sessions, err := db.Q.MxSessionsGet(ctx, projectID)
	if err != nil {
		err = fmt.Errorf("failed to get sessions from db: %w", err)
		return nil, err
	}

	var res []types.Exec

	for _, s := range sessions {
		s := s
		if !listTerminated && s.Status == db.UnweaveSessionStatusTerminated {
			continue
		}
		connInfo := &ConnectionInfoV1{}
		if err := json.Unmarshal(s.ConnectionInfo, connInfo); err != nil {
			return nil, fmt.Errorf("failed to unmarshal connection info: %w", err)
		}
		session := types.Exec{
			ID:   s.ID,
			Name: s.Name,
			SSHKey: types.SSHKey{
				Name:      s.SshKeyName,
				PublicKey: &s.PublicKey,
				CreatedAt: &s.SshKeyCreatedAt,
			},
			Connection: &types.ConnectionInfo{
				Host: connInfo.Host,
				Port: connInfo.Port,
				User: connInfo.User,
			},
			Status:     types.NodeStatus(s.Status),
			CreatedAt:  &s.CreatedAt,
			NodeTypeID: s.NodeID,
			Region:     s.Region,
			Provider:   types.Provider(s.Provider),
		}
		res = append(res, session)
	}
	return res, nil
}

func (s *ExecService) Watch(ctx context.Context, sessionID string) error {
	session, err := db.Q.MxSessionGet(ctx, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &types.Error{
				Code:    http.StatusNotFound,
				Message: "Session not found",
			}
		}
		return fmt.Errorf("failed to get session from db: %w", err)
	}

	rt, err := s.srv.InitializeRuntime(ctx, types.Provider(session.Provider))
	if err != nil {
		return fmt.Errorf("failed to initialize runtime: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	statusch, errch := rt.Node.Watch(ctx, session.NodeID)

	log.Ctx(ctx).Info().Msgf("Starting to watch session %s", sessionID)

	go func() {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			case status := <-statusch:
				log.Ctx(ctx).
					Info().
					Str(SessionStatusCtxKey, string(status)).
					Msg("session status changed")

				if status == types.StatusRunning {
					if e := updateConnectionInfo(ctx, rt.Node, session.NodeID, sessionID); e != nil {
						// We mark the error in the DB but don't terminate the node. This
						// is left to the user to do manually. Perhaps this should be
						// changed in the future but for now, it might help debugging.
						handleSessionError(ctx, sessionID, e, "Failed to update connection info")
						// TODO: we should perhaps do some retries here
						return
					}
				}

				params := db.SessionStatusUpdateParams{
					ID:     sessionID,
					Status: db.UnweaveSessionStatus(status),
					ReadyAt: sql.NullTime{
						Time:  time.Now(),
						Valid: true,
					},
				}
				if e := db.Q.SessionStatusUpdate(ctx, params); e != nil {
					log.Ctx(ctx).Error().Err(e).Msg("failed to update session status")
					return
				}
				if status == types.StatusTerminated {
					return
				}
			case e := <-errch:
				log.Ctx(ctx).Error().Err(e).Msg("Error while watching session")

				// This means we failed to watch the session. This should ideally never
				// happen. Since we don't know the cause of this error, let's play it safe
				// and terminate the node. This will mean the user will lose their work
				// but the alternative is to have a runaway node that drains all their credit.
				if err := s.Terminate(ctx, sessionID); err != nil {
					log.Ctx(ctx).Error().Err(err).Msg("failed to terminate session on failure to watch")
				}
				handleSessionError(ctx, sessionID, e, "Failed to watch session")
				return
			}
		}
	}()

	return nil
}

func (s *ExecService) Terminate(ctx context.Context, sessionID string) error {
	sess, err := db.Q.MxSessionGet(ctx, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &types.Error{
				Code:       http.StatusNotFound,
				Message:    "Session not found",
				Suggestion: "Make sure the session id is valid",
			}
		}
		return fmt.Errorf("failed to fetch session from db %q: %w", sessionID, err)
	}

	provider := types.Provider(sess.Provider)
	rt, err := s.srv.InitializeRuntime(ctx, provider)
	if err != nil {
		return fmt.Errorf("failed to create runtime %q: %w", sess.Provider, err)
	}

	ctx = log.With().
		Stringer(types.RuntimeProviderKey, s.srv.runtime.Node.GetProvider()).
		Logger().
		WithContext(ctx)

	if err = rt.Node.TerminateNode(ctx, sess.NodeID); err != nil {
		return fmt.Errorf("failed to terminate node: %w", err)
	}
	if err = s.srv.vault.DeleteSecret(ctx, sess.NodeID); err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("Failed to delete secret for node %q", sess.NodeID)
	}

	params := db.SessionStatusUpdateParams{
		ID:     sessionID,
		Status: db.UnweaveSessionStatusTerminated,
		ExitedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	if err = db.Q.SessionStatusUpdate(ctx, params); err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Msgf("Failed to set session %q as terminated", sessionID)
	}
	return nil
}
