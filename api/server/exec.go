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

func (c ConnectionInfoV1) GetConnectionInfo() *types.ConnectionInfo {
	return &types.ConnectionInfo{
		Host: c.Host,
		Port: c.Port,
		User: c.User,
	}
}

type NodeMetadataV1 struct {
	ID             string           `json:"id"`
	TypeID         string           `json:"typeID"`
	Price          int              `json:"price"`
	VCPUs          int              `json:"vcpus"`
	Memory         int              `json:"memory"`
	HDD            int              `json:"hdd"`
	GpuType        string           `json:"gpuType"`
	GPUCount       int              `json:"gpuCount"`
	GPUMemory      int              `json:"gpuMemory"`
	ConnectionInfo ConnectionInfoV1 `json:"connection_info"`
}

func (m NodeMetadataV1) GetHardwareSpec() types.HardwareSpec {
	return types.HardwareSpec{
		GPU: types.GPU{
			Count: types.HardwareRequestRange{
				Min: m.GPUCount,
				Max: m.GPUCount,
			},
			Type: m.GpuType,
			RAM: types.HardwareRequestRange{
				Min: m.GPUMemory,
				Max: m.GPUMemory,
			},
		},
		CPU: types.HardwareRequestRange{
			Min: m.VCPUs,
			Max: m.VCPUs,
		},
		RAM: types.HardwareRequestRange{
			Min: m.Memory,
			Max: m.Memory,
		},
		HDD: types.HardwareRequestRange{
			Min: m.HDD,
			Max: m.HDD,
		},
	}
}

func DBNodeMetadataFromNode(node types.Node) NodeMetadataV1 {
	n := NodeMetadataV1{
		ID:        node.ID,
		TypeID:    node.TypeID,
		Price:     node.Price,
		VCPUs:     node.Specs.CPU.Min,
		Memory:    node.Specs.RAM.Min,
		HDD:       node.Specs.HDD.Min,
		GpuType:   node.Specs.GPU.Type,
		GPUCount:  node.Specs.GPU.Count.Min,
		GPUMemory: node.Specs.GPU.RAM.Min,

		ConnectionInfo: ConnectionInfoV1{
			Version: 1,
			Host:    node.Host,
			Port:    node.Port,
			User:    node.User,
		},
	}
	return n
}

func handleExecError(execID string, err error, msg string) {
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

	params := db.ExecSetErrorParams{
		ID: execID,
		Error: sql.NullString{
			String: msg,
			Valid:  true,
		},
	}
	if err = db.Q.ExecSetError(ctx, params); err != nil {
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

	// Public key wasn't provided and key wasn't found by name
	if sshPublicKey == nil {
		return types.SSHKey{}, &types.Error{
			Code:    http.StatusBadRequest,
			Message: "SSH key not found",
		}
	}
	if sshKeyName == nil || *sshKeyName == "" {
		sshKeyName = tools.Stringy("uw:" + random.GenerateRandomPhrase(4, "-") + ".pub")
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
	params := db.ExecUpdateConnectionInfoParams{
		ID:             execID,
		ConnectionInfo: connInfoJSON,
	}
	if e := db.Q.ExecUpdateConnectionInfo(ctx, params); e != nil {
		return fmt.Errorf("failed to update connection info: %w", e)
	}
	return nil
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

	keys, err := s.setupUserCreds(ctx, params.SSHKeyName, params.SSHPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to setup credentials: %w", err)
	}

	node, err := s.assignNode(ctx, types.SetSpecDefaultValues(params.HardwareSpec), params.Region, keys)
	if err != nil {
		return nil, fmt.Errorf("failed to assign node: %w", err)
	}

	if _, err := s.srv.vault.SetSecret(ctx, *keys[0].PrivateKey, &node.ID); err != nil {
		return nil, fmt.Errorf("failed to store private key: %w", err)
	}

	buildID, imageURI, err := s.getExecImage(ctx, projectID, params.Image)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image for exec: %w", err)
	}

	var exec *types.Exec

	execCfg := types.ExecConfig{
		Image:   imageURI,
		Command: params.Command,
		Keys:    keys[1:], // Only mount user keys into exec
		Volumes: nil,      // TODO: implement attaching volumes
		Src:     params.Source,
	}

	gitCfg := types.GitConfig{
		CommitID: params.CommitID,
		GitURL:   params.GitURL,
	}

	exec, err = s.init(ctx, projectID, node, execCfg, gitCfg, buildID, imageURI)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize exec: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get image uri: %w", err)
	}

	go func() {
		c := context.Background()
		c = log.With().
			Str(UserIDCtxKey, s.srv.cid).
			Str(ProjectIDCtxKey, projectID).
			Str(ExecIDCtxKey, exec.ID).
			Logger().WithContext(c)

		if e := s.srv.Exec.Watch(c, exec.ID); e != nil {
			log.Ctx(c).Error().Err(e).Msgf("Failed to watch exec")
		}
	}()

	return exec, nil
}

func (s *ExecService) Get(ctx context.Context, execID string) (*types.Exec, error) {
	dbs, err := db.Q.MxExecGet(ctx, execID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &types.Error{
				Code:    http.StatusNotFound,
				Message: "Exec not found",
			}
		}
		return nil, fmt.Errorf("failed to get session from db: %w", err)
	}

	metadata := &NodeMetadataV1{}
	if err := json.Unmarshal(dbs.Metadata, metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal connection info: %w", err)
	}

	session := types.NewExec(
		execID,
		dbs.Name,
		types.SSHKey{
			Name:      dbs.SshKeyName,
			PublicKey: &dbs.PublicKey,
			CreatedAt: &dbs.SshKeyCreatedAt,
		},
		"",
		metadata.ConnectionInfo.GetConnectionInfo(),
		types.Status(dbs.Status),
		&dbs.CreatedAt,
		metadata.TypeID,
		dbs.Region,
		types.Provider(dbs.Provider),
		metadata.GetHardwareSpec(),
		false, // Assuming `hasPersistentFS` value is always `false` in this case
	)
	return session, nil
}

func (s *ExecService) List(ctx context.Context, projectID string, listAll bool) ([]types.Exec, error) {
	sessions, err := db.Q.MxExecsGet(ctx, projectID)
	if err != nil {
		err = fmt.Errorf("failed to get sessions from db: %w", err)
		return nil, err
	}

	var res []types.Exec

	for _, s := range sessions {
		s := s
		if !listAll && (s.Status == db.UnweaveExecStatusTerminated || s.Status == db.UnweaveExecStatusError) {
			continue
		}
		metadata := &NodeMetadataV1{}
		if err := json.Unmarshal(s.Metadata, metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal connection info: %w", err)
		}
		session := types.NewExec(
			s.ID,
			s.Name,
			types.SSHKey{
				Name:      s.SshKeyName,
				PublicKey: &s.PublicKey,
				CreatedAt: &s.SshKeyCreatedAt,
			},
			"",
			metadata.ConnectionInfo.GetConnectionInfo(),
			types.Status(s.Status),
			&s.CreatedAt,
			metadata.TypeID,
			s.Region,
			types.Provider(s.Provider),
			metadata.GetHardwareSpec(),
			false, // Assuming `hasPersistentFS` value is always `false` in this case
		)
		res = append(res, *session)
	}
	return res, nil
}

func (s *ExecService) Watch(ctx context.Context, execID string) error {
	exec, err := db.Q.MxExecGet(ctx, execID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &types.Error{
				Code:    http.StatusNotFound,
				Message: "Exec not found",
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
					Str(ExecStatusCtxKey, string(status)).
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
						handleExecError(execID, e, "Failed to update connection info")
						// TODO: we should perhaps do some retries here
						return
					}
				}

				params := db.ExecStatusUpdateParams{
					ID:     execID,
					Status: db.UnweaveExecStatus(status),
					ReadyAt: sql.NullTime{
						Time:  time.Now(),
						Valid: true,
					},
				}
				if e := db.Q.ExecStatusUpdate(ctx, params); e != nil {
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
				handleExecError(execID, e, "Failed to watch exec")
				return
			}
		}
	}()

	return nil
}

func (s *ExecService) Terminate(ctx context.Context, execID string) error {
	exec, err := db.Q.MxExecGet(ctx, execID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &types.Error{
				Code:       http.StatusNotFound,
				Message:    "Exec not found",
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

	if err = rt.Exec.Terminate(ctx, exec.ID); err != nil {
		return fmt.Errorf("failed to terminate node: %w", err)
	}
	if err = s.srv.vault.DeleteSecret(ctx, exec.NodeID); err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("Failed to delete secret for node %q", exec.NodeID)
	}

	params := db.ExecStatusUpdateParams{
		ID:     execID,
		Status: db.UnweaveExecStatusTerminated,
		ExitedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	if err = db.Q.ExecStatusUpdate(ctx, params); err != nil {
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
