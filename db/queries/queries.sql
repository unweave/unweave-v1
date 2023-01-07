-- name: ProjectCreate :one
insert into unweave.projects (name, owner_id)
values ($1, $2)
returning id;

-- name: ProjectGet :one
select *
from unweave.projects
where id = $1;

-- name: SessionCreate :one
insert into unweave.sessions (node_id, created_by, project_id, runtime, ssh_key_id)
values ($1, $2, $3, $4, (select id
                         from unweave.ssh_keys
                         where name = @ssh_key_name
                           and owner_id = $2))
returning id;

-- name: SessionGet :one
select *
from unweave.sessions
where id = $1;

-- name: SessionAndProjectGet :one
select s.id,
       s.node_id,
       s.created_by,
       s.created_at,
       s.ready_at,
       s.exited_at,
       s.status,
       s.runtime,
       s.ssh_key_id,
       p.id         as project_id,
       p.name       as project_name,
       p.owner_id   as project_owner_id,
       p.created_at as project_created_at
from unweave.sessions s
         join unweave.projects p on p.id = s.project_id
where s.id = $1;

-- name: SessionsGet :many
select sessions.id, ssh_keys.name as ssh_key_name, sessions.status
from unweave.sessions
         left join unweave.ssh_keys
                   on ssh_keys.id = sessions.ssh_key_id
where project_id = $1
order by unweave.sessions.created_at desc
limit $2 offset $3;

-- name: SessionSetTerminated :exec
update unweave.sessions
set status = unweave.session_status('terminated')
where id = $1;

-- name: SSHKeyAdd :exec
insert INTO unweave.ssh_keys (owner_id, name, public_key)
values ($1, $2, $3);

-- name: SSHKeysGet :many
select *
from unweave.ssh_keys
where owner_id = $1;

-- name: SSHKeyGetByName :one
select *
from unweave.ssh_keys
where name = $1;

-- name: SSHKeyGetByPublicKey :one
select *
from unweave.ssh_keys
where public_key = $1;
