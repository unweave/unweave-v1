-- name: BuildCreate :one
insert into unweave.build (project_id, builder_type)
values ($1, $2)
returning id;

-- name: BuildGet :one
select *
from unweave.build
where id = $1;

-- name: BuildUpdate :exec
update unweave.build
set status    = $2,
    meta_data = $3
where id = $1;

-- name: ProjectGet :one
select *
from unweave.project
where id = $1;

-- name: SessionCreate :one
insert into unweave.session (node_id, created_by, project_id, provider, ssh_key_id,
                             region, name, connection_info)
values ($1, $2, $3, $4, (select id
                         from unweave.ssh_key as ssh_keys
                         where ssh_keys.name = @ssh_key_name
                           and owner_id = $2), $5, $6, $7)
returning id;

-- name: SessionGet :one
select *
from unweave.session
where id = $1;

-- name: SessionGetAllActive :many
select *
from unweave.session
where status = 'initializing'
   or status = 'running';

-- name: SessionUpdateConnectionInfo :exec
update unweave.session
set connection_info = $2
where id = $1;

-- name: SessionsGet :many
select session.id, ssh_key.name as ssh_key_name, session.status
from unweave.session
         left join unweave.ssh_key
                   on ssh_key.id = session.ssh_key_id
where project_id = $1
order by unweave.session.created_at desc
limit $2 offset $3;

-- name: SessionSetError :exec
update unweave.session
set status = 'error'::unweave.session_status,
    error  = $2
where id = $1;

-- name: SessionStatusUpdate :exec
update unweave.session
set status = $2
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

-- name: MxSessionGet :one
select s.id,
       s.status,
       s.node_id,
       s.provider,
       s.region,
       s.created_at,
       s.connection_info,
       ssh_key.name       as ssh_key_name,
       ssh_key.public_key,
       ssh_key.created_at as ssh_key_created_at
from unweave.session as s
         join unweave.ssh_key on s.ssh_key_id = ssh_key.id
where s.id = $1;

-- name: MxSessionsGet :many
select s.id,
       s.status,
       s.node_id,
       s.provider,
       s.region,
       s.created_at,
       s.connection_info,
       ssh_key.name       as ssh_key_name,
       ssh_key.public_key,
       ssh_key.created_at as ssh_key_created_at
from unweave.session as s
         join unweave.ssh_key on s.ssh_key_id = ssh_key.id
where s.project_id = $1;