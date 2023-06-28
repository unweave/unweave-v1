package execsrv

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/tools"
)

type postgresStore struct{}

func NewPostgresStore() Store {
	return postgresStore{}
}

func (p postgresStore) Create(projectID string, exec types.Exec) error {
	ctx := context.Background()

	if projectID == "" {
		return fmt.Errorf("an Exec must be attached to a project")
	}

	// Every exec should be created with a public key
	publicKeys := filterNullPublicKeys(exec.Keys)
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
		return fmt.Errorf("failed to create SSH keys: %w", err)
	}

	params := db.ExecCreateParams{
		ID:           exec.ID,
		CreatedBy:    exec.CreatedBy,
		ProjectID:    projectID,
		Region:       exec.Region,
		Name:         exec.Name,
		Spec:         spec,
		Metadata:     metadata,
		BuildID:      db.NullStringFrom(exec.BuildID),
		Image:        exec.Image,
		Provider:     exec.Provider.String(),
		CommitID:     db.NullStringFrom(exec.CommitID),
		GitRemoteUrl: db.NullStringFrom(exec.GitURL),
		Command:      []string{},
	}

	if err = db.Q.ExecCreate(ctx, params); err != nil {
		return fmt.Errorf("failed to create exec: %w", err)
	}

	keys, err := getSSHKeysByPublicKey(ctx, exec.CreatedBy, extractPublicKeys(publicKeys))
	if err != nil {
		return err
	}

	err = addSSHKeyToExec(ctx, exec, keys)
	if err != nil {
		return fmt.Errorf("failed to add SSH key to exec: %w", err)
	}

	for _, volume := range exec.Volumes {
		err := db.Q.ExecVolumeCreate(context.Background(), db.ExecVolumeCreateParams{
			ExecID:    exec.ID,
			VolumeID:  volume.VolumeID,
			MountPath: volume.MountPath,
		})
		if err != nil {
			return fmt.Errorf("failed to assign volumes to exec with error: %w", err)
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
	keyRefs, err := db.Q.ExecSSHKeysGetByExecID(ctx, id)
	if err != nil {
		return types.Exec{}, err
	}
	keyIDs := tools.MapToStrings(keyRefs, func(key db.UnweaveExecSshKey) string {
		return key.SshKeyID
	})

	keys, err := db.Q.SSHKeysGetByIDs(ctx, keyIDs)
	if err != nil {
		return types.Exec{}, err
	}
	volumes, err := db.Q.ExecVolumeGet(ctx, id)

	return dbExecToExec(exec, dbExecVolumesToVolumes(volumes), dbSSHKeyToSSHKey(keys)), nil
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

	for idx, exec := range execs {
		keyRefs, err := db.Q.ExecSSHKeysGetByExecID(ctx, exec.ID)
		if err != nil {
			return nil, err
		}

		keyIDs := tools.MapToStrings(keyRefs, func(key db.UnweaveExecSshKey) string {
			return key.SshKeyID
		})

		keys, err := db.Q.SSHKeysGetByIDs(ctx, keyIDs)
		if err != nil {
			return []types.Exec{}, err
		}

		volumes, err := db.Q.ExecVolumeGet(ctx, exec.ID)
		if err != nil {
			return nil, err
		}

		res[idx] = dbExecToExec(exec, dbExecVolumesToVolumes(volumes), dbSSHKeyToSSHKey(keys))
	}

	return res, nil
}

func (p postgresStore) Delete(id string) error {
	// Execs should be soft deleted
	err := db.Q.ExecVolumeDelete(context.Background(), id)
	if err != nil {
		return fmt.Errorf("failed to unassign volumes for exec with error: %w", err)
	}

	if err = p.UpdateStatus(id, types.StatusTerminated, time.Time{}, time.Now()); err != nil {
		return fmt.Errorf("failed to update exec status in store: %w", err)
	}

	return nil
}

func (p postgresStore) Update(id string, exec types.Exec) error {
	//TODO implement me
	panic("implement me")
}

// UpdateStatus updates exec status and relevant timestamps
func (p postgresStore) UpdateStatus(id string, status types.Status, setReadyAt, setExitedAt time.Time) error {
	params := db.ExecStatusUpdateParams{
		ID:       id,
		Status:   db.UnweaveExecStatus(status),
		ReadyAt:  db.NullTimeFrom(setReadyAt),
		ExitedAt: db.NullTimeFrom(setExitedAt),
	}
	if e := db.Q.ExecStatusUpdate(context.Background(), params); e != nil {
		return fmt.Errorf("failed to update exec status: %w", e)
	}
	return nil
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

func createSSHKeys(ctx context.Context, createdByID string, keys []types.SSHKey) error {
	for _, key := range keys {
		key := key
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

func dbSSHKeyToSSHKey(keys []db.UnweaveSshKey) []types.SSHKey {
	res := make([]types.SSHKey, len(keys))

	for idx, k := range keys {
		k := k
		res[idx] = types.SSHKey{
			Name:       k.Name,
			PublicKey:  &k.PublicKey,
			PrivateKey: nil,
			CreatedAt:  &k.CreatedAt,
		}
	}

	return res
}

func dbExecToExec(dbe db.UnweaveExec, volumes []types.ExecVolume, keys []types.SSHKey) types.Exec {
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
		Volumes:   volumes,
		Network:   metadataFromJSON.GetExecNetwork(),
		Spec:      *spec,
		CommitID:  commitID,
		GitURL:    githubRemoteURL,
		Region:    dbe.Region,
		Provider:  types.Provider(dbe.Provider),
	}
}

func dbExecVolumesToVolumes(volumes []db.UnweaveExecVolume) []types.ExecVolume {
	out := make([]types.ExecVolume, len(volumes))

	for i, volume := range volumes {
		out[i] = dbExecVolumeToExecVolume(volume)
	}

	return out
}

func dbExecVolumeToExecVolume(volume db.UnweaveExecVolume) types.ExecVolume {
	return types.ExecVolume{
		VolumeID:  volume.VolumeID,
		MountPath: volume.MountPath,
	}
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

func filterNullPublicKeys(keys []types.SSHKey) []types.SSHKey {
	filteredKeys := make([]types.SSHKey, 0)

	for _, key := range keys {
		if key.PublicKey != nil {
			filteredKeys = append(filteredKeys, key)
		}
	}

	return filteredKeys
}

func extractPublicKeys(keys []types.SSHKey) []string {
	filteredStrings := make([]string, 0)

	for _, key := range keys {
		if key.PublicKey != nil {
			filteredStrings = append(filteredStrings, *key.PublicKey)
		}
	}

	return filteredStrings
}
