-- name: ExecCreate :exec
insert into unweave.exec (id, created_by, project_id,
                          region, name, spec, metadata, commit_id, git_remote_url,
                          command,
                          build_id, image, provider)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);

-- name: ExecGet :one
select *
from unweave.exec
where id = @id_or_name
   or name = @id_or_name;

-- name: ExecGetAllActive :many
select *
from unweave.exec
where status = 'initializing'
   or status = 'running';

-- name: ExecListByProvider :many
select *
from unweave.exec as e
where e.provider = $1;

-- name: ExecList :many
select *
from unweave.exec as e
where (e.provider = coalesce(sqlc.narg('filter_provider'), e.provider))
  and project_id = coalesce(sqlc.narg('filter_project_id'), project_id)
  and ((@filter_active = true and (status = 'pending' or status = 'initializing' or status = 'running'))
    or @filter_active = false);

-- name: ExecListActiveByProvider :many
select *
from unweave.exec as e
where provider = $1
  and (status = 'initializing'
    or status = 'running');

-- name: ExecUpdateConnectionInfo :exec
update unweave.exec
set metadata = jsonb_set(metadata, '{connection_info}', @connection_info::jsonb)
where id = $1;

-- name: ExecUpdateNetwork :exec
update unweave.exec
set metadata = jsonb_set(metadata, '{http_service}', @http_service::jsonb)
where id = $1;

-- name: ExecSetError :exec
update unweave.exec
set status = 'error'::unweave.exec_status,
    error  = $2
where id = $1;

-- name: ExecSetFailed :exec
update unweave.exec
set status = 'failed'::unweave.exec_status,
    error  = $2
where id = $1;

-- name: ExecStatusUpdate :exec
update unweave.exec
set status    = $2,
    ready_at  = coalesce(sqlc.narg('ready_at'), ready_at),
    exited_at = coalesce(sqlc.narg('exited_at'), exited_at)
where id = $1;
