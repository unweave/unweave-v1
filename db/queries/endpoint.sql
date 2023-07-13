-- name: EndpointCreate :exec
INSERT INTO unweave.endpoint (id, name, project_id, http_address, created_at) VALUES ($1, $2, $3, $4, $5);

-- name: EndpointGet :one
SELECT id, name, project_id, http_address, created_at, deleted_at
FROM unweave.endpoint
WHERE id = $1 OR (name = $1 AND project_id = $2);

-- name: EndpointDelete :exec
DELETE FROM unweave.endpoint WHERE id = $1;

-- name: EndpointsForProject :many
SELECT id, name, project_id, http_address, created_at, deleted_at
FROM unweave.endpoint
WHERE project_id = $1;

-- name: EndpointEval :many
SELECT endpoint_id, eval_id FROM unweave.endpoint_eval WHERE endpoint_id = $1;

-- name: EndpointEvalAttach :exec
INSERT INTO unweave.endpoint_eval (endpoint_id, eval_id) VALUES ($1, $2);

-- name: EndpointCheckCreate :exec
INSERT INTO unweave.endpoint_check (id, endpoint_id, project_id) VALUES ($1, $2, $3);

-- name: EndpointCheck :one
SELECT id, endpoint_id, project_id, created_at FROM unweave.endpoint_check WHERE id = $1;

-- name: EndpointCheckStepCreate :exec
INSERT INTO unweave.endpoint_check_step (id, check_id, eval_id, input) VALUES ($1, $2, $3, $4);

-- name: EndpointCheckStepUpdate :exec
UPDATE unweave.endpoint_check_step
SET input = coalesce(sqlc.narg('input'), input),
    output = coalesce(sqlc.narg('output'), output),
    assertion = coalesce(sqlc.narg('assertion'), assertion)
WHERE id = sqlc.narg('id');

-- name: EndpointCheckSteps :many
SELECT id, check_id, eval_id, input, output, assertion
FROM unweave.endpoint_check_step
WHERE check_id = $1;


