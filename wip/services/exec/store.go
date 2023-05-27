package exec

import (
	"context"
	"database/sql"

	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/db"
	"github.com/unweave/unweave/tools/random"
)

type postgresStore struct{}

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

	dbp := db.ExecCreateParams{
		ID:           exec.ID,
		NodeID:       "",
		CreatedBy:    exec.CreatedBy,
		ProjectID:    project,
		Region:       exec.Region,
		Name:         random.GenerateRandomPhrase(4, "-"),
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

func (p postgresStore) Get(project, id string) (types.Exec, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresStore) List(project string) ([]types.Exec, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresStore) ListAll() ([]types.Exec, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresStore) Delete(project, id string) error {
	//TODO implement me
	panic("implement me")
}

func (p postgresStore) Update(project, id string, exec types.Exec) error {
	//TODO implement me
	panic("implement me")
}
