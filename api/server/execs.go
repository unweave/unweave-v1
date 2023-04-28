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

var DefaultImageURI = "ubuntu:latest"

type ConnectionInfoV1 struct {
	Version int    `json:"version"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	User    string `json:"user"`
}

type NodeMetadataV1 struct {
	ID             string           `json:"id"`
	Price          int              `json:"price"`
	VCPUs          int              `json:"vcpus"`
	Memory         int              `json:"memory"`
	Storage        int              `json:"storage"`
	GpuType        string           `json:"gpuType"`
	GPUCount       int              `json:"gpuCount"`
	GPUMemory      int              `json:"gpuMemory"`
	ConnectionInfo ConnectionInfoV1 `json:"connection_info"`
}

func DBNodeMetadataFromNode(node types.Node) NodeMetadataV1 {
	gpuMem := 0
	if node.Specs.GPUMemory != nil {
		gpuMem = *node.Specs.GPUMemory
	}
	n := NodeMetadataV1{
		ID:        node.ID,
		Price:     node.Price,
		VCPUs:     node.Specs.VCPUs,
		Memory:    node.Specs.Memory,
		Storage:   node.Specs.Storage,
		GpuType:   node.Specs.GPUType,
		GPUCount:  node.Specs.GPUCount,
		GPUMemory: gpuMem,

		ConnectionInfo: ConnectionInfoV1{
			Version: 1,
			Host:    node.Host,
			Port:    node.Port,
			User:    node.User,
		},
	}
	return n
}

func handleSessionError(execID string, err error, msg string) {
	// Make sure this doesn't fail because of a parent cancelled context
	ctx := context.Background()
	ctx = log.With().Logger().WithContext(ctx)
	ctx = log.Ctx(ctx).With().Str(ExecIDCtxKey, execID).Logger().WithContext(ctx)

	var e *types.Error
	if errors.As(err, &e) {
		msg += ": " + e.Message
	}
	msg += ": " + err.Error()

	log.Ctx(ctx).Error().Err(err).Msg(msg)

	params := db.SessionSetErrorParams{
		ID: execID,
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

func updateConnectionInfo(rt runtime.Exec, nodeID string, execID string) error {
	// New ctx to make sure this doesn't fail because of a parent cancelled context
	ctx := context.Background()
	connInfo, err := rt.GetConnectionInfo(ctx, execID)
	if err != nil {
		return fmt.Errorf("failed to get connection info: %w", err)
	}

	log.Info().Str("node_id", nodeID).Str("exec_id", execID).Msgf("Updating connection info %v", connInfo)
	connInfoJSON, err := json.Marshal(connInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal connection info: %w", err)
	}
	params := db.SessionUpdateConnectionInfoParams{
		ID:             execID,
		ConnectionInfo: connInfoJSON,
	}
	if e := db.Q.SessionUpdateConnectionInfo(ctx, params); e != nil {
		return fmt.Errorf("failed to update connection info: %w", e)
	}
	return nil
}

func updateExecStatus(ctx context.Context, execID string, status types.Status) {
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

func (s *ExecService) Create(ctx context.Context, projectID string, params types.ExecCreateParams, filesystemID *string) (*types.Exec, error) {
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
	node.OwnerID = s.srv.aid

	metadata := DBNodeMetadataFromNode(node)
	metadataJSON, err := json.Marshal(&metadata)
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
		Metadata:  metadataJSON,
		Status:    string(types.StatusInitializing),
		OwnerID:   s.srv.aid,
		SshKeyIds: []string{sshKey[0].ID},
	}
	if err = db.Q.NodeCreate(ctx, np); err != nil {
		return nil, fmt.Errorf("failed to create node in db: %w", err)
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

	bid := sql.NullString{}
	imageURI := DefaultImageURI
	buildID := params.Ctx.BuildID

	if buildID == nil {
		project, err := db.Q.ProjectGet(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("failed to get project: %w", err)
		}
		if project.DefaultBuildID.Valid && project.DefaultBuildID.String != "" {
			buildID = &project.DefaultBuildID.String
		}
	}
	if buildID != nil && *buildID != "" {
		imageURI, err = s.srv.Builder.GetImageURI(ctx, *params.Ctx.BuildID)
		if err != nil {
			log.Error().Err(err).Msgf("failed to get image uri: %w", err)
			return nil, fmt.Errorf("failed to get image uri: %w", err)
		}
		bid = sql.NullString{String: *buildID, Valid: true}
	}

	execID, err := rt.Exec.Init(ctx, node, []types.SSHKey{userKey}, imageURI, filesystemID)
	if err != nil {
		go handleSessionError(execID, err, "failed to init session")
		return nil, fmt.Errorf("failed to init session: %w", err)
	}
	dbp := db.SessionCreateParams{
		ID:           execID,
		NodeID:       node.ID,
		CreatedBy:    s.srv.cid,
		ProjectID:    projectID,
		Region:       node.Region,
		Name:         random.GenerateRandomPhrase(4, "-"),
		Metadata:     metadataJSON,
		CommitID:     commitID,
		GitRemoteUrl: gitRemoteURL,
		Command:      command,
		BuildID:      bid,
		PersistFs:    filesystemID != nil,
		SshKeyName:   userKey.Name,
	}

	if err := db.Q.SessionCreate(ctx, dbp); err != nil {
		return nil, fmt.Errorf("failed to create session in db: %w", err)
	}
	ctx = log.With().Str(ExecIDCtxKey, execID).Logger().WithContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get image uri: %w", err)
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
	}

	go func() {
		c := context.Background()
		c = log.With().
			Str(UserIDCtxKey, s.srv.cid).
			Str(ProjectIDCtxKey, projectID).
			Str(ExecIDCtxKey, session.ID).
			Logger().WithContext(c)

		if e := s.srv.Exec.Watch(c, session.ID); e != nil {
			log.Ctx(ctx).Error().Err(e).Msgf("Failed to watch session")
		}
	}()

	return session, nil
}

func (s *ExecService) CreateFromSnapshot(ctx context.Context, projectID string, filesystemID string) error {
	// TODO: if filesystem is not a root filesystem, return error
	return nil
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

	metadata := &NodeMetadataV1{}
	if err := json.Unmarshal(dbs.Metadata, metadata); err != nil {
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
			Host: metadata.ConnectionInfo.Host,
			Port: metadata.ConnectionInfo.Port,
			User: metadata.ConnectionInfo.User,
		},
		Status:     types.Status(dbs.Status),
		CreatedAt:  &dbs.CreatedAt,
		NodeTypeID: dbs.NodeID,
		Region:     dbs.Region,
		Provider:   types.Provider(dbs.Provider),
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
		metadata := &NodeMetadataV1{}
		if err := json.Unmarshal(s.Metadata, metadata); err != nil {
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
				Host: metadata.ConnectionInfo.Host,
				Port: metadata.ConnectionInfo.Port,
				User: metadata.ConnectionInfo.User,
			},
			Status:     types.Status(s.Status),
			CreatedAt:  &s.CreatedAt,
			NodeTypeID: s.NodeID,
			Region:     s.Region,
			Provider:   types.Provider(s.Provider),
		}
		res = append(res, session)
	}
	return res, nil
}

func (s *ExecService) Watch(ctx context.Context, execID string) error {
	exec, err := db.Q.MxSessionGet(ctx, execID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &types.Error{
				Code:    http.StatusNotFound,
				Message: "Session not found",
			}
		}
		return fmt.Errorf("failed to get exec from db: %w", err)
	}

	rt, err := s.srv.InitializeRuntime(ctx, types.Provider(exec.Provider))
	if err != nil {
		return fmt.Errorf("failed to initialize runtime: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	statusch, errch := rt.Exec.Watch(ctx, exec.ID)

	log.Ctx(ctx).Info().Msgf("Starting to watch exec %s", execID)

	go func() {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				log.Ctx(ctx).Info().Msgf("Ctx Done. Stopping to watch exec %s", execID)
				return
			case status := <-statusch:
				log.Ctx(ctx).
					Info().
					Str(SessionStatusCtxKey, string(status)).
					Msg("Exec status changed")

				if status == types.StatusRunning {
					if e := updateConnectionInfo(rt.Exec, exec.NodeID, execID); e != nil {
						// Use new context to make sure terminate is not cancelled
						terminateCtx := context.Background()
						terminateCtx = log.With().Logger().WithContext(terminateCtx)

						if terminateErr := s.Terminate(terminateCtx, execID); terminateErr != nil {
							log.Error().
								Err(terminateErr).
								Msgf("failed to terminate exec %q on failure to watch", execID)
						}
						handleSessionError(execID, e, "Failed to update connection info")
						// TODO: we should perhaps do some retries here
						return
					}
				}

				params := db.SessionStatusUpdateParams{
					ID:     execID,
					Status: db.UnweaveSessionStatus(status),
					ReadyAt: sql.NullTime{
						Time:  time.Now(),
						Valid: true,
					},
				}
				if e := db.Q.SessionStatusUpdate(ctx, params); e != nil {
					log.Ctx(ctx).Error().Err(e).Msg("failed to update exec status")
					return
				}

				if status == types.StatusTerminated {
					log.Ctx(ctx).Info().Msgf("Exec %q exited", execID)
					// Use new context to make sure terminate is not cancelled
					terminateCtx := context.Background()
					terminateCtx = log.With().Logger().WithContext(terminateCtx)

					// Clean up before returning. This will should be a no-op if the pod was
					// already deleted. This is particularly going to happen when a pod is
					// naturally terminated at end of exec.
					if err = s.Terminate(terminateCtx, execID); err != nil {
						log.Warn().Err(err).Msgf("Failed to terminate exec %q", execID)
					}
					return
				}
			case e := <-errch:
				log.Ctx(ctx).Error().Err(e).Msg("Error while watching exec")

				// This means we failed to watch the exec. This should ideally never
				// happen. Since we don't know the cause of this error, let's play it safe
				// and terminate the node. This will mean the user will lose their work
				// but the alternative is to have a runaway node that drains all their credit.
				if err := s.Terminate(ctx, execID); err != nil {
					log.Ctx(ctx).Error().Err(err).Msg("failed to terminate exec on failure to watch")
				}
				handleSessionError(execID, e, "Failed to watch exec")
				return
			}
		}
	}()

	return nil
}

func (s *ExecService) Snapshot(ctx context.Context, execID string) error {
	exec, err := db.Q.MxSessionGet(ctx, execID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &types.Error{
				Code:       http.StatusNotFound,
				Message:    "Session not found",
				Suggestion: "Make sure the session id is valid",
			}
		}
		return fmt.Errorf("failed to fetch session from db %q: %w", execID, err)
	}

	if !exec.PersistFs {
		log.Ctx(ctx).Info().Msgf("Exec %q is not configured to persist filesystem. No-op.", execID)
		return nil
	}

	fs, e := db.Q.FilesystemGetByExecID(ctx, execID)
	if e != nil {
		return fmt.Errorf("failed to get filesystem for exec %q: %w", execID, e)
	}

	provider := types.Provider(exec.Provider)
	rt, err := s.srv.InitializeRuntime(ctx, provider)
	if err != nil {
		return fmt.Errorf("failed to create runtime %q: %w", exec.Provider, err)
	}

	err = rt.Exec.SnapshotFS(ctx, execID, fs.ID)
	if err != nil {
		return fmt.Errorf("failed to snapshot filesystem: %w", err)
	}

	return nil
}

func (s *ExecService) Terminate(ctx context.Context, execID string) error {
	exec, err := db.Q.MxSessionGet(ctx, execID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &types.Error{
				Code:       http.StatusNotFound,
				Message:    "Session not found",
				Suggestion: "Make sure the session id is valid",
			}
		}
		return fmt.Errorf("failed to fetch session from db %q: %w", execID, err)
	}

	if string(exec.Status) == string(types.StatusTerminated) {
		log.Ctx(ctx).Info().Msgf("Exec %q is already terminated. No-op.", execID)
		return nil
	}

	// Use new context to make sure we terminate the pod even if the context is cancelled
	ctx = context.Background()
	ctx = log.With().Str("execID", execID).Logger().WithContext(ctx)

	provider := types.Provider(exec.Provider)
	rt, err := s.srv.InitializeRuntime(ctx, provider)
	if err != nil {
		return fmt.Errorf("failed to create runtime %q: %w", exec.Provider, err)
	}

	ctx = log.With().
		Stringer(types.RuntimeProviderKey, s.srv.runtime.Node.GetProvider()).
		Logger().
		WithContext(ctx)

	if err = s.Snapshot(ctx, execID); err != nil {
		// We don't want to fail the termination because of this. Log and continue.
		log.Ctx(ctx).Error().Err(err).Msgf("Failed to snapshot filesystem for exec %q", execID)
	}
	if err = rt.Exec.Terminate(ctx, exec.ID); err != nil {
		return fmt.Errorf("failed to terminate node: %w", err)
	}
	if err = s.srv.vault.DeleteSecret(ctx, exec.NodeID); err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("Failed to delete secret for node %q", exec.NodeID)
	}

	params := db.SessionStatusUpdateParams{
		ID:     execID,
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
			Msgf("Failed to set session %q as terminated", execID)
	}

	np := db.NodeStatusUpdateParams{
		ID:           exec.NodeID,
		Status:       "terminated",
		TerminatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}
	if err = db.Q.NodeStatusUpdate(ctx, np); err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Msgf("Failed to set node %q as terminated", exec.NodeID)
	}
	return nil
}
