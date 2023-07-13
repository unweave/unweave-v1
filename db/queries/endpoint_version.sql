-- name: EndpointVersion :one
SELECT id, endpoint_id, exec_id, project_id, http_address, primary_version, created_at, deleted_at FROM unweave.endpoint_version WHERE id = $1;

-- name: EndpointVersionList :many
SELECT id, endpoint_id, exec_id, project_id, http_address, primary_version, created_at, deleted_at FROM unweave.endpoint_version WHERE endpoint_id = $1;

-- name: EndpointVersionCreate :exec
INSERT INTO unweave.endpoint_version (id, endpoint_id, exec_id, project_id, http_address, created_at) VALUES ($1, $2, $3, $4, $5, $6);

-- name: EndpointVersionDemote :exec
UPDATE unweave.endpoint_version
SET primary_version = FALSE
WHERE endpoint_id = $1;

-- name: EndpointVersionPromote :exec
UPDATE unweave.endpoint_version
SET primary_version = TRUE
WHERE id = $1;
