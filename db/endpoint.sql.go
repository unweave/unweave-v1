// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: endpoint.sql

package db

import (
	"context"
)

const EndpointCreate = `-- name: EndpointCreate :exec
INSERT INTO unweave.endpoint (id, exec_id, project_id) VALUES ($1, $2, $3)
`

type EndpointCreateParams struct {
	ID        string `json:"id"`
	ExecID    string `json:"execID"`
	ProjectID string `json:"projectID"`
}

func (q *Queries) EndpointCreate(ctx context.Context, arg EndpointCreateParams) error {
	_, err := q.db.ExecContext(ctx, EndpointCreate, arg.ID, arg.ExecID, arg.ProjectID)
	return err
}

const EndpointDelete = `-- name: EndpointDelete :exec
DELETE FROM unweave.endpoint WHERE id = $1
`

func (q *Queries) EndpointDelete(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, EndpointDelete, id)
	return err
}

const EndpointEval = `-- name: EndpointEval :many
SELECT endpoint_id, eval_id FROM unweave.endpoint_eval WHERE endpoint_id = $1
`

func (q *Queries) EndpointEval(ctx context.Context, endpointID string) ([]UnweaveEndpointEval, error) {
	rows, err := q.db.QueryContext(ctx, EndpointEval, endpointID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UnweaveEndpointEval
	for rows.Next() {
		var i UnweaveEndpointEval
		if err := rows.Scan(&i.EndpointID, &i.EvalID); err != nil {
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

const EndpointEvalAttach = `-- name: EndpointEvalAttach :exec
INSERT INTO unweave.endpoint_eval (endpoint_id, eval_id) VALUES ($1, $2)
`

type EndpointEvalAttachParams struct {
	EndpointID string `json:"endpointID"`
	EvalID     string `json:"evalID"`
}

func (q *Queries) EndpointEvalAttach(ctx context.Context, arg EndpointEvalAttachParams) error {
	_, err := q.db.ExecContext(ctx, EndpointEvalAttach, arg.EndpointID, arg.EvalID)
	return err
}

const EndpointGet = `-- name: EndpointGet :one
SELECT id, exec_id, project_id, created_at, deleted_at FROM unweave.endpoint WHERE id = $1
`

func (q *Queries) EndpointGet(ctx context.Context, id string) (UnweaveEndpoint, error) {
	row := q.db.QueryRowContext(ctx, EndpointGet, id)
	var i UnweaveEndpoint
	err := row.Scan(
		&i.ID,
		&i.ExecID,
		&i.ProjectID,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const EndpointsForProject = `-- name: EndpointsForProject :many
SELECT id, exec_id, project_id, created_at, deleted_at FROM unweave.endpoint WHERE project_id = $1
`

func (q *Queries) EndpointsForProject(ctx context.Context, projectID string) ([]UnweaveEndpoint, error) {
	rows, err := q.db.QueryContext(ctx, EndpointsForProject, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UnweaveEndpoint
	for rows.Next() {
		var i UnweaveEndpoint
		if err := rows.Scan(
			&i.ID,
			&i.ExecID,
			&i.ProjectID,
			&i.CreatedAt,
			&i.DeletedAt,
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
