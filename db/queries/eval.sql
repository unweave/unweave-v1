-- name: EvalList :many
SELECT id, exec_id, http_address, created_at FROM unweave.eval WHERE id = ANY($1::text[]);

-- name: EvalCreate :exec
INSERT INTO unweave.eval (id, exec_id, http_address, project_id) VALUES ($1, $2, $3, $4);

-- name: EvalDelete :exec
DELETE FROM unweave.eval WHERE id = $1;

-- name: EvalGet :one
SELECT id, exec_id, http_address, project_id FROM unweave.eval WHERE id = $1;

-- name: EvalListForProject :many
SELECT id, exec_id, http_address, project_id from unweave.eval WHERE project_id = $1;

