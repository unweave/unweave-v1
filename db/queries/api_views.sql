-------------------------------------------------------------------
-- The queries below return data in the format expected by the API.
-------------------------------------------------------------------

-- name: MxExecGet :one
select e.id,
       e.name,
       e.status,
       e.provider,
       e.region,
       e.created_at,
       e.metadata
from unweave.exec as e
where e.id = $1;

-- name: MxExecsGet :many
select e.id,
       e.name,
       e.status,
       e.provider,
       e.region,
       e.created_at,
       e.metadata
from unweave.exec as e
where e.project_id = $1;
