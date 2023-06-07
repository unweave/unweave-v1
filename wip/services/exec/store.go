package exec

import (
	"context"
	"database/sql"
	"fmt"

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
		Provider:     exec.Provider.String(),
		SshKeyName:   "",
	}

	if err := db.Q.ExecCreate(context.Background(), dbp); err != nil {
		// TODO: parse db errors: project not found, ssh key not found,
		return err
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

	return dbExecToExec(exec), nil
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

func (p postgresStore) List(project string) ([]types.Exec, error) {
	ctx := context.Background()

	params := db.ExecsGetParams{
		ProjectID: project,
		Limit:     1000,
		Offset:    0,
	}
	execs, err := db.Q.ExecsGet(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	res := make([]types.Exec, len(execs))

	for idx, exec := range execs {
		exec := exec
		res[idx] = dbExecToExec(exec)
	}
	return res, nil
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
		execs[idx] = dbExecToExec(dbe)
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

func dbExecToExec(dbe db.UnweaveExec) types.Exec {
	var bid *string
	if dbe.BuildID.Valid {
		bid = &dbe.BuildID.String
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
