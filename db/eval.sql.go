// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: eval.sql

package db

import (
	"context"
	"time"

	"github.com/lib/pq"
)

const EvalCreate = `-- name: EvalCreate :exec
INSERT INTO unweave.eval (id, exec_id, project_id) VALUES ($1, $2, $3)
`

type EvalCreateParams struct {
	ID        string `json:"id"`
	ExecID    string `json:"execID"`
	ProjectID string `json:"projectID"`
}

func (q *Queries) EvalCreate(ctx context.Context, arg EvalCreateParams) error {
	_, err := q.db.ExecContext(ctx, EvalCreate, arg.ID, arg.ExecID, arg.ProjectID)
	return err
}

const EvalDelete = `-- name: EvalDelete :exec
DELETE FROM unweave.eval WHERE id = $1
`

func (q *Queries) EvalDelete(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, EvalDelete, id)
	return err
}

const EvalGet = `-- name: EvalGet :one
SELECT id, exec_id, project_id FROM unweave.eval WHERE id = $1
`

type EvalGetRow struct {
	ID        string `json:"id"`
	ExecID    string `json:"execID"`
	ProjectID string `json:"projectID"`
}

func (q *Queries) EvalGet(ctx context.Context, id string) (EvalGetRow, error) {
	row := q.db.QueryRowContext(ctx, EvalGet, id)
	var i EvalGetRow
	err := row.Scan(&i.ID, &i.ExecID, &i.ProjectID)
	return i, err
}

const EvalList = `-- name: EvalList :many
SELECT id, exec_id, created_at FROM unweave.eval WHERE id = ANY($1::text[])
`

type EvalListRow struct {
	ID        string    `json:"id"`
	ExecID    string    `json:"execID"`
	CreatedAt time.Time `json:"createdAt"`
}

func (q *Queries) EvalList(ctx context.Context, dollar_1 []string) ([]EvalListRow, error) {
	rows, err := q.db.QueryContext(ctx, EvalList, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EvalListRow
	for rows.Next() {
		var i EvalListRow
		if err := rows.Scan(&i.ID, &i.ExecID, &i.CreatedAt); err != nil {
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
