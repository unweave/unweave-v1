-- name: ExecCreate :exec
insert into unweave.exec (id, node_id, created_by, project_id, ssh_key_id,
                          region, name, metadata, commit_id, git_remote_url, command,
                          build_id, image)
values ($1, $2, $3, $4, (select id
                         from unweave.ssh_key as ssh_keys
                         where ssh_keys.name = @ssh_key_name
                           and owner_id = $3), $5, $6, $7, $8, $9, $10, $11, $12);

-- name: ExecGet :one
select *
from unweave.exec
where id = @id_or_name or name = @id_or_name;

-- name: ExecGetAllActive :many
select *
from unweave.exec
where status = 'initializing'
   or status = 'running';

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
