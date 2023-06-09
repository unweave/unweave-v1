package exec

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
)

type postgresStore struct{}

func NewPostgresStore() Store {
	return postgresStore{}
}

func (p postgresStore) Create(ctx context.Context, project string, exec types.Exec) error {
	if project == "" {
		return fmt.Errorf("an Exec must be attached to a project")
	}

	// every exec should be created with a public key
	publicKeys := types.FilterKeysWithPublicKey(exec.Keys)
	if len(publicKeys) > 1 || len(publicKeys) == 0 {
		return fmt.Errorf("an Exec must be created with one and only one SSH public key")
	}
	if exec.Name == "" {
		return fmt.Errorf("an Exec must be named")
	}

	spec, err := json.Marshal(&exec.Spec)
	if err != nil {
		return fmt.Errorf("failed to marshal spec to JSON: %w", err)
	}

	metadata, err := json.Marshal(&types.NodeMetadataV1{})
	if err != nil {
		return fmt.Errorf("failed to marshal metadata to JSON: %w", err)
	}

	err = createSSHKeys(ctx, exec.CreatedBy, publicKeys)
	if err != nil {
		return err
	}
	err = createExec(project, spec, metadata, exec)
	if err != nil {
		return err
	}
	keys, err := getSSHKeysByPublicKey(ctx, exec.CreatedBy, types.GetPublicKeys(publicKeys))
	if err != nil {
		return err
	}
	err = addSSHKeyToExec(ctx, exec, keys)
	if err != nil {
		return err
	}

	return nil
}

func createExec(projectID string, spec []byte, metadata []byte, exec types.Exec) error {
	if err := db.Q.ExecCreate(context.Background(), db.ExecCreateParams{
		ID:        exec.ID,
		CreatedBy: exec.CreatedBy,
		ProjectID: projectID,
		Region:    exec.Region,
		Name:      exec.Name,
		Spec:      spec,
		Metadata:  metadata,
		BuildID:   db.NullStringFrom(exec.BuildID),
		Image:     exec.Image,
		Provider:  exec.Provider.String(),

		// Note: These fields are members of the Exec, but currently unused in any feature.
		CommitID:     db.NullStringFrom(exec.CommitID),
		GitRemoteUrl: db.NullStringFrom(exec.GitURL),
		Command:      []string{},
	}); err != nil {
		return err
	}

	return nil
}

func createSSHKeys(ctx context.Context, createdByID string, keys []types.SSHKey) error {
	for _, key := range keys {
		exists, err := db.Q.SSHKeyGetByPublicKey(ctx, db.SSHKeyGetByPublicKeyParams{
			PublicKey: *key.PublicKey,
			OwnerID:   createdByID,
		})
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if exists.PublicKey != "" {
			continue
		}
		err = db.Q.SSHKeyAdd(ctx, db.SSHKeyAddParams{
			OwnerID:   createdByID,
			Name:      key.Name,
			PublicKey: *key.PublicKey,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func getSSHKeysByPublicKey(ctx context.Context, ownerID string, pubs []string) ([]db.UnweaveSshKey, error) {
	keys := make([]db.UnweaveSshKey, 0, len(pubs))

	for _, pub := range pubs {
		key, err := db.Q.SSHKeyGetByPublicKey(ctx, db.SSHKeyGetByPublicKeyParams{
			PublicKey: pub,
			OwnerID:   ownerID,
		})
		if err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}

	return keys, nil
}

func addSSHKeyToExec(ctx context.Context, exec types.Exec, keys []db.UnweaveSshKey) error {
	for _, key := range keys {
		err := db.Q.ExecSSHKeyInsert(ctx, db.ExecSSHKeyInsertParams{
			ExecID:   exec.ID,
			SshKeyID: key.ID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (p postgresStore) Get(id string) (types.Exec, error) {
	ctx := context.Background()

	exec, err := db.Q.ExecGet(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Exec{}, ErrNotFound
		}
		return types.Exec{}, err
	}
	keyRefs, err := db.Q.ExecSSHKeysByExecIDGet(ctx, id)
	if err != nil {
		return types.Exec{}, err
	}
	keyIDs := db.MapStrings(keyRefs, func(key db.UnweaveExecSshKey) string {
		return key.SshKeyID
	})

	keys, err := db.Q.SSHKeysGetByIDs(ctx, keyIDs)
	if err != nil {
		return types.Exec{}, err
	}

	return dbExecToExec(exec, dbSSHKeyToSSHKey(keys)), nil
}

func (p postgresStore) GetDriver(id string) (string, error) {
	exec, err := p.Get(id)
	if err != nil {
		return "", err
	}
	// TODO: this should eventually be changed to a driver instead of a provider when
	//  we have more than one driver
	return exec.Provider.String(), nil
}

func (p postgresStore) List(projectID *string, filterProvider *types.Provider, filterActive bool) ([]types.Exec, error) {
	ctx := context.Background()

	var project, provider sql.NullString

	if projectID != nil {
		project = sql.NullString{String: *projectID, Valid: true}
	}
	if filterProvider != nil {
		provider = sql.NullString{String: filterProvider.String(), Valid: true}
	}

	params := db.ExecListParams{
		FilterProvider:  provider,
		FilterProjectID: project,
		FilterActive:    filterActive,
	}
	execs, err := db.Q.ExecList(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list execs: %w", err)
	}
	res := make([]types.Exec, len(execs))

	for _, exec := range execs {
		keyRefs, err := db.Q.ExecSSHKeysByExecIDGet(ctx, exec.ID)
		if err != nil {
			return nil, err
		}
		keyIDs := db.MapStrings(keyRefs, func(key db.UnweaveExecSshKey) string {
			return key.SshKeyID
		})
		keys, err := db.Q.SSHKeysGetByIDs(ctx, keyIDs)
		if err != nil {
			return []types.Exec{}, err
		}

		res = append(res, dbExecToExec(exec, dbSSHKeyToSSHKey(keys)))
	}

	return res, nil
}

func (p postgresStore) Delete(project, id string) error {
	//TODO implement me
	panic("implement me")
}

func (p postgresStore) Update(id string, exec types.Exec) error {
	//TODO implement me
	panic("implement me")
}

func (p postgresStore) UpdateStatus(id string, status types.Status) error {
	params := db.ExecStatusUpdateParams{
		ID:     id,
		Status: db.UnweaveExecStatus(status),
	}
	if e := db.Q.ExecStatusUpdate(context.Background(), params); e != nil {
		return fmt.Errorf("failed to update exec status: %w", e)
	}
	return nil
}

func dbSSHKeyToSSHKey(ks []db.UnweaveSshKey) (res []types.SSHKey) {
	for _, k := range ks {
		res = append(res, types.SSHKey{
			Name:       k.Name,
			PublicKey:  &k.PublicKey,
			PrivateKey: nil,
			CreatedAt:  &k.CreatedAt,
		})
	}

	return res
}

func dbExecToExec(dbe db.UnweaveExec, keys []types.SSHKey) types.Exec {
	var bid *string
	if dbe.BuildID.Valid {
		bid = &dbe.BuildID.String
	}
	var commitID *string
	if dbe.CommitID.Valid {
		commitID = &dbe.CommitID.String
	}
	var githubRemoteURL *string
	if dbe.GitRemoteUrl.Valid {
		githubRemoteURL = &dbe.GitRemoteUrl.String
	}

	metadataFromJSON, err := types.NodeMetadataFromJSON(dbe.Metadata)
	if err != nil {
		log.Err(err).Msg("failed to properly unmarshal node metadata, metadata will not be parsed")
	}

	spec, err := types.HardwareSpecFromJSON(dbe.Spec)
	if err != nil {
		log.Err(err).Msg("failed to properly unmarshal exec spec, spec will not be parsed")
	}
	if spec == nil {
		spec = new(types.HardwareSpec)
	}

	return types.Exec{
		ID:        dbe.ID,
		Name:      dbe.Name,
		CreatedAt: dbe.CreatedAt,
		CreatedBy: dbe.CreatedBy,
		Image:     dbe.Image,
		BuildID:   bid,
		Status:    types.Status(dbe.Status),
		Command:   dbe.Command,
		Keys:      keys,
		Volumes:   nil,
		Network:   metadataFromJSON.GetExecNetwork(),
		Spec:      *spec,
		CommitID:  commitID,
		GitURL:    githubRemoteURL,
		Region:    dbe.Region,
		Provider:  types.Provider(dbe.Provider),
	}
}
