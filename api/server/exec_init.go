package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/tools/random"
)

func (s *ExecService) assignNode(ctx context.Context, nodeTypeID string, region *string, NodeGPUCount int, keys []types.SSHKey) (types.Node, error) {
	owner := s.srv.aid
	user := s.srv.cid

	node, err := s.srv.runtime.Node.InitNode(ctx, keys, nodeTypeID, region, NodeGPUCount)
	if err != nil {
		return types.Node{}, fmt.Errorf("failed to init node: %w", err)
	}
	node.OwnerID = owner

	metadata := DBNodeMetadataFromNode(node)
	metadataJSON, err := json.Marshal(&metadata)
	if err != nil {
		return types.Node{}, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	sshKey, err := db.Q.SSHKeysGet(ctx, user)
	if err != nil {
		return types.Node{}, fmt.Errorf("failed to get ssh key ids: %w", err)
	}

	np := db.NodeCreateParams{
		ID:        node.ID,
		Provider:  string(s.srv.runtime.Node.GetProvider()),
		Region:    node.Region,
		Metadata:  metadataJSON,
		Status:    string(types.StatusInitializing),
		OwnerID:   owner,
		SshKeyIds: []string{sshKey[0].ID},
	}
	if err = db.Q.NodeCreate(ctx, np); err != nil {
		return types.Node{}, fmt.Errorf("failed to create node in db: %w", err)
	}

	return node, nil
}

func (s *ExecService) getExecImage(ctx context.Context, projectID string, imageOrBuild *string) (buildID string, imageURI string, err error) {
	imageURI = DefaultImageURI

	// No image or build specified, get default image for project
	if imageOrBuild == nil {
		// Get default image for project
		project, err := db.Q.ProjectGet(ctx, projectID)
		if err != nil {
			return "", "", fmt.Errorf("failed to get project: %w", err)
		}

		if !project.DefaultBuildID.Valid || project.DefaultBuildID.String == "" {
			return "", imageURI, nil
		}

		imageOrBuild = &project.DefaultBuildID.String
	}

	if *imageOrBuild != "" {
		build, err := db.Q.BuildGet(ctx, *imageOrBuild)
		if err != nil && err != sql.ErrNoRows {
			return "", "", fmt.Errorf("failed to get build: %w", err)
		}

		// Must be a referencing a build
		if err == nil {
			imageURI, err = s.srv.Builder.GetImageURI(ctx, build.ID)
			if err != nil {
				return "", "", fmt.Errorf("failed to get image uri: %w", err)
			}
			return build.ID, imageURI, nil
		}

		// Must be a public image
		return "", *imageOrBuild, nil
	}

	// No default image for project and no image specified, use default
	return "", imageURI, nil
}

func (s *ExecService) setupUserCreds(ctx context.Context, keyName, pubKey *string) ([]types.SSHKey, error) {
	user := s.srv.cid
	userKey, err := fetchCredentials(ctx, user, keyName, pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to setup credentials: %w", err)
	}
	prv, pub, err := generateSSHKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ssh key pair: %w", err)
	}

	adminKey := types.SSHKey{
		Name:       "umk-" + random.GenerateRandomAdjectiveNounTriplet(),
		PublicKey:  &pub,
		PrivateKey: &prv,
		CreatedAt:  nil,
	}
	keys := []types.SSHKey{adminKey, userKey}

	if err = registerCredentials(ctx, s.srv.runtime, keys); err != nil {
		return nil, fmt.Errorf("failed to register credentials: %w", err)
	}

	return keys, nil
}

func (s *ExecService) init(ctx context.Context, projectID string, node types.Node, cfg types.ExecConfig, gitCfg types.GitConfig, buildID, imageURI string) (*types.Exec, error) {
	var command []string
	var commitID, gitRemoteURL sql.NullString

	if cfg.Command != nil {
		command = cfg.Command
	}
	if gitCfg.CommitID != nil {
		commitID = sql.NullString{String: *gitCfg.CommitID, Valid: true}
	}
	if gitCfg.GitURL != nil {
		gitRemoteURL = sql.NullString{String: *gitCfg.GitURL, Valid: true}
	}

	execID, err := s.srv.runtime.Exec.Init(ctx, node, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init exec: %w", err)
	}
	metadata := DBNodeMetadataFromNode(node)
	metadataJSON, err := json.Marshal(&metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}
	createdAt := time.Now()

	bid := sql.NullString{}
	if buildID != "" {
		bid = sql.NullString{String: buildID, Valid: true}
	}

	dbp := db.ExecCreateParams{
		ID:           execID,
		NodeID:       node.ID,
		CreatedBy:    s.srv.cid,
		ProjectID:    projectID,
		Region:       node.Region,
		Name:         random.GenerateRandomPhrase(4, "-"),
		Metadata:     metadataJSON, // This is currently the same as the node metadata. Will change in the future.
		CommitID:     commitID,
		GitRemoteUrl: gitRemoteURL,
		Command:      command,
		BuildID:      bid,
		Image:        imageURI,
		PersistFs:    len(cfg.Volumes) != 0, // TODO: implement this properly
		SshKeyName:   cfg.Keys[0].Name,      // TODO: support multiple keys
	}

	if err := db.Q.ExecCreate(ctx, dbp); err != nil {
		return nil, fmt.Errorf("failed to create exec in db: %w", err)
	}

	exec := &types.Exec{
		ID:           execID,
		Name:         dbp.Name,
		SSHKey:       cfg.Keys[0],
		Image:        cfg.Image,
		Connection:   nil,
		Status:       types.StatusInitializing,
		CreatedAt:    &createdAt,
		NodeTypeID:   node.TypeID,
		Region:       node.Region,
		Provider:     node.Provider,
		PersistentFS: false,
	}
	return exec, nil
}
