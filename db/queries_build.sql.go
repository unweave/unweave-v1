// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: queries_build.sql

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

const BuildCreate = `-- name: BuildCreate :one
insert into unweave.build (project_id, builder_type, name, created_by, started_at)
values ($1, $2, $3, $4, case
                            when $5::timestamptz = '0001-01-01 00:00:00 UTC'::timestamptz
                                then now()
                            else $5::timestamptz end)
    returning id
`

type BuildCreateParams struct {
	ProjectID   string    `json:"projectID"`
	BuilderType string    `json:"builderType"`
	Name        string    `json:"name"`
	CreatedBy   string    `json:"createdBy"`
	StartedAt   time.Time `json:"startedAt"`
}

func (q *Queries) BuildCreate(ctx context.Context, arg BuildCreateParams) (string, error) {
	row := q.db.QueryRowContext(ctx, BuildCreate,
		arg.ProjectID,
		arg.BuilderType,
		arg.Name,
		arg.CreatedBy,
		arg.StartedAt,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

const BuildGet = `-- name: BuildGet :one
select id, name, project_id, builder_type, status, created_by, created_at, started_at, finished_at, updated_at, meta_data
from unweave.build
where id = $1
`

func (q *Queries) BuildGet(ctx context.Context, id string) (UnweaveBuild, error) {
	row := q.db.QueryRowContext(ctx, BuildGet, id)
	var i UnweaveBuild
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ProjectID,
		&i.BuilderType,
		&i.Status,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.StartedAt,
		&i.FinishedAt,
		&i.UpdatedAt,
		&i.MetaData,
	)
	return i, err
}

const BuildGetUsedBy = `-- name: BuildGetUsedBy :many
select s.id, s.name, s.region, s.created_by, s.created_at, s.ready_at, s.exited_at, s.status, s.project_id, s.error, s.build_id, s.spec, s.commit_id, s.git_remote_url, s.command, s.metadata, s.image, s.provider, n.provider
from (select id from unweave.build as ub where ub.id = $1) as b
         join unweave.exec s
              on s.build_id = b.id
         join unweave.node as n on s.node_id = n.id
`

type BuildGetUsedByRow struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Region       string            `json:"region"`
	CreatedBy    string            `json:"createdBy"`
	CreatedAt    time.Time         `json:"createdAt"`
	ReadyAt      sql.NullTime      `json:"readyAt"`
	ExitedAt     sql.NullTime      `json:"exitedAt"`
	Status       UnweaveExecStatus `json:"status"`
	ProjectID    string            `json:"projectID"`
	Error        sql.NullString    `json:"error"`
	BuildID      sql.NullString    `json:"buildID"`
	Spec         json.RawMessage   `json:"spec"`
	CommitID     sql.NullString    `json:"commitID"`
	GitRemoteUrl sql.NullString    `json:"gitRemoteUrl"`
	Command      []string          `json:"command"`
	Metadata     json.RawMessage   `json:"metadata"`
	Image        string            `json:"image"`
	Provider     string            `json:"provider"`
	Provider_2   string            `json:"provider2"`
}

func (q *Queries) BuildGetUsedBy(ctx context.Context, id string) ([]BuildGetUsedByRow, error) {
	rows, err := q.db.QueryContext(ctx, BuildGetUsedBy, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BuildGetUsedByRow
	for rows.Next() {
		var i BuildGetUsedByRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Region,
			&i.CreatedBy,
			&i.CreatedAt,
			&i.ReadyAt,
			&i.ExitedAt,
			&i.Status,
			&i.ProjectID,
			&i.Error,
			&i.BuildID,
			&i.Spec,
			&i.CommitID,
			&i.GitRemoteUrl,
			pq.Array(&i.Command),
			&i.Metadata,
			&i.Image,
			&i.Provider,
			&i.Provider_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const BuildUpdate = `-- name: BuildUpdate :exec
update unweave.build
set status      = $2,
    meta_data   = $3,
    started_at  = coalesce(
            nullif($4::timestamptz, '0001-01-01 00:00:00 UTC'::timestamptz),
            started_at),
    finished_at = coalesce(
            nullif($5::timestamptz, '0001-01-01 00:00:00 UTC'::timestamptz),
            finished_at)
where id = $1
`

type BuildUpdateParams struct {
	ID         string             `json:"id"`
	Status     UnweaveBuildStatus `json:"status"`
	MetaData   json.RawMessage    `json:"metaData"`
	StartedAt  time.Time          `json:"startedAt"`
	FinishedAt time.Time          `json:"finishedAt"`
}

func (q *Queries) BuildUpdate(ctx context.Context, arg BuildUpdateParams) error {
	_, err := q.db.ExecContext(ctx, BuildUpdate,
		arg.ID,
		arg.Status,
		arg.MetaData,
		arg.StartedAt,
		arg.FinishedAt,
	)
	return err
}
