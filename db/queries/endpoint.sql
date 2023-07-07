-- name: EndpointCreate :exec
INSERT INTO unweave.endpoint (id, exec_id, project_id) VALUES ($1, $2, $3);

-- name: EndpointGet :one
SELECT id, exec_id, project_id, created_at, deleted_at FROM unweave.endpoint WHERE id = $1;

-- name: EndpointDelete :exec
DELETE FROM unweave.endpoint WHERE id = $1;

-- name: EndpointEval :many
SELECT endpoint_id, eval_id FROM unweave.endpoint_eval WHERE endpoint_id = $1;

-- name: EndpointEvalAttach :exec
INSERT INTO unweave.endpoint_eval (endpoint_id, eval_id) VALUES ($1, $2);

-- name: EndpointsForProject :many
SELECT id, exec_id, project_id, created_at, deleted_at FROM unweave.endpoint WHERE project_id = $1;
