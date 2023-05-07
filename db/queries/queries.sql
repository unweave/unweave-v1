-- name: BuildCreate :one
insert into unweave.build (project_id, builder_type, name, created_by, started_at)
values ($1, $2, $3, $4, case
                            when @started_at::timestamptz = '0001-01-01 00:00:00 UTC'::timestamptz
                                then now()
                            else @started_at::timestamptz end)
returning id;


-- name: BuildGet :one
select *
from unweave.build
where id = $1;

-- name: BuildGetUsedBy :many
select s.*, n.provider
from (select id from unweave.build as ub where ub.id = $1) as b
         join unweave.exec s
              on s.build_id = b.id
         join unweave.node as n on s.node_id = n.id;

-- name: BuildUpdate :exec
update unweave.build
set status      = $2,
    meta_data   = $3,
    started_at  = coalesce(
            nullif(@started_at::timestamptz, '0001-01-01 00:00:00 UTC'::timestamptz),
            started_at),
    finished_at = coalesce(
            nullif(@finished_at::timestamptz, '0001-01-01 00:00:00 UTC'::timestamptz),
            finished_at)
where id = $1;


-- name: FilesystemCreate :one
insert into unweave.filesystem (name, project_id, owner_id, exec_id, src_path)
values ($1, $2, $3, $4, $5)
returning *;

-- name: FilesystemCreateVersion :one
select *
from unweave.insert_filesystem_version(@filesystem_id, @exec_id);

-- name: FilesystemGet :one
select *
from unweave.filesystem
where id = $1;

-- name: FilesystemGetLatestVersion :one
select *
from unweave.filesystem_version
where filesystem_id = $1
order by version desc
limit 1;

-- name: FilesystemGetByProject :one
select *
from unweave.filesystem
where project_id = @project_id
  and (name = @name_or_id or id = @name_or_id);

-- name: FilesystemGetByExecID :one
select b.*
from (select filesystem_id
      from unweave.filesystem_version
      where filesystem_version.exec_id = $1) as bv
         join unweave.filesystem b on b.id = filesystem_id;

-- name: FilesystemVersionGet :one
select *
from unweave.filesystem_version
where exec_id = $1;

-- name: FilesystemVersionAddBuildID :exec
update unweave.filesystem_version
set build_id = $2
where exec_id = $1;


-- name: ProjectGet :one
select *
from unweave.project
where id = $1;

-- name: NodeCreate :exec
select unweave.insert_node(
               @id,
               @provider,
               @region,
               @metadata :: jsonb,
               @status,
               @owner_id,
               @ssh_key_ids :: text[]
           );

-- name: NodeStatusUpdate :exec
update unweave.node
set status        = $2,
    ready_at      = coalesce(sqlc.narg('ready_at'), ready_at),
    terminated_at = coalesce(sqlc.narg('terminated_at'), terminated_at)
where id = $1;


-- name: ExecCreate :exec
insert into unweave.exec (id, node_id, created_by, project_id, ssh_key_id,
                             region, name, metadata, commit_id, git_remote_url, command,
                             build_id, image,  persist_fs)
values ($1, $2, $3, $4, (select id
                         from unweave.ssh_key as ssh_keys
                         where ssh_keys.name = @ssh_key_name
                           and owner_id = $3), $5, $6, $7, $8, $9, $10, $11, $12, $13);

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
select exec.id, ssh_key.name as ssh_key_name, exec.status
from unweave.exec
         left join unweave.ssh_key
                   on ssh_key.id = exec.ssh_key_id
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

-- name: SSHKeyAdd :exec
insert into unweave.ssh_key (owner_id, name, public_key)
values ($1, $2, $3);

-- name: SSHKeysGet :many
select *
from unweave.ssh_key
where owner_id = $1;

-- name: SSHKeyGetByName :one
select *
from unweave.ssh_key
where name = $1
  and owner_id = $2;

-- name: SSHKeyGetByPublicKey :one
select *
from unweave.ssh_key
where public_key = $1
  and owner_id = $2;


-------------------------------------------------------------------
-- The queries below return data in the format expected by the API.
-------------------------------------------------------------------

-- name: MxExecGet :one
select s.id,
       s.name,
       s.status,
       s.node_id,
       n.provider,
       s.region,
       s.created_at,
       s.metadata,
       s.persist_fs,
       ssh_key.name       as ssh_key_name,
       ssh_key.public_key,
       ssh_key.created_at as ssh_key_created_at
from unweave.exec as s
         join unweave.ssh_key on s.ssh_key_id = ssh_key.id
         join unweave.node as n on s.node_id = n.id
where s.id = $1;

-- name: MxExecsGet :many
select s.id,
       s.name,
       s.status,
       s.node_id,
       n.provider,
       s.region,
       s.created_at,
       s.metadata,
       s.persist_fs,
       ssh_key.name       as ssh_key_name,
       ssh_key.public_key,
       ssh_key.created_at as ssh_key_created_at
from unweave.exec as s
         join unweave.ssh_key on s.ssh_key_id = ssh_key.id
         join unweave.node as n on s.node_id = n.id
where s.project_id = $1;