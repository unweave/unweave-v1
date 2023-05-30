package exec

import (
	"context"
	"database/sql"

	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/tools/random"
)

type postgresStore struct{}

func NewPostgresStore() Store {
	return &postgresStore{}
}

func (p postgresStore) Create(project string, exec types.Exec) error {
	var command []string
	var commitID, gitRemoteURL sql.NullString

	if exec.Command != nil {
		command = exec.Command
	}
	if exec.CommitID != nil {
		commitID = sql.NullString{String: *exec.CommitID, Valid: true}
	}
	if exec.GitURL != nil {
		gitRemoteURL = sql.NullString{String: *exec.GitURL, Valid: true}
	}
	bid := sql.NullString{}
	if exec.BuildID != nil {
		bid = sql.NullString{String: *exec.BuildID, Valid: true}
	}
	if exec.Name != "" {
		exec.Name = random.GenerateRandomPhrase(4, "-")
	}

	dbp := db.ExecCreateParams{
		ID:           exec.ID,
		NodeID:       "",
		CreatedBy:    exec.CreatedBy,
		ProjectID:    project,
		Region:       exec.Region,
		Name:         exec.Name,
		Metadata:     nil,
		CommitID:     commitID,
		GitRemoteUrl: gitRemoteURL,
		Command:      command,
		BuildID:      bid,
		Image:        exec.Image,
		SshKeyName:   "",
	}

	if err := db.Q.ExecCreate(context.Background(), dbp); err != nil {
		// TODO: parse db errors: project not found, ssh key not found,
		return err
	}
	return nil
}

func (p postgresStore) Get(id string) (types.Exec, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresStore) GetDriver(id string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresStore) List(project string) ([]types.Exec, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresStore) ListByProvider(provider types.Provider, filterActive bool) ([]types.Exec, error) {
	ctx := context.Background()

	var err error
	var dbExecs []db.UnweaveExec

	if filterActive {
		dbExecs, err = db.Q.ExecListActiveByProvider(ctx, provider.String())
		if err != nil {
			return nil, err
		}
	} else {
		dbExecs, err = db.Q.ExecListByProvider(ctx, provider.String())
		if err != nil {
			return nil, err
		}
	}

	execs := make([]types.Exec, len(dbExecs))

	for idx, dbe := range dbExecs {
		dbe := dbe
		var bid *string
		if dbe.BuildID.Valid {
			bid = &dbe.BuildID.String
		}

		execs[idx] = types.Exec{
			ID:        dbe.ID,
			Name:      dbe.Name,
			CreatedAt: dbe.CreatedAt,
			CreatedBy: dbe.CreatedBy,
			Image:     dbe.Image,
			BuildID:   bid,
			Status:    types.Status(dbe.Status),
			Command:   dbe.Command,
			Keys:      nil,
			Volumes:   nil,
			Network:   types.ExecNetwork{},
			Spec:      types.HardwareSpec{},
			CommitID:  nil,
			GitURL:    nil,
			Region:    "",
			Provider:  types.Provider(dbe.Provider),
		}
	}
	return execs, nil
}

func (p postgresStore) Delete(project, id string) error {
	//TODO implement me
	panic("implement me")
}

func (p postgresStore) Update(id string, exec types.Exec) error {
	//TODO implement me
	panic("implement me")
}
