-- name: ExecCreate :exec
insert into unweave.exec (id, created_by, project_id,
                          region, name, spec, metadata, commit_id, git_remote_url, command,
                          build_id, image, provider)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,  $13);

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

-- name: ExecsGet :many
select *
from unweave.exec
where project_id = $1
order by unweave.exec.created_at desc
limit $2 offset $3;

-- name: ExecSetError :exec
update unweave.exec
set status = 'error'::unweave.exec_status,
    error  = $2
where id = $1;

-- name: ExecStatusUpdate :exec
update unweave.exec
set status    = $2,
    ready_at  = coalesce(sqlc.narg('ready_at'), ready_at),
    exited_at = coalesce(sqlc.narg('exited_at'), exited_at)
where id = $1;
