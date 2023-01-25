-- name: ProjectCreate :one
insert into unweave.project (name, owner_id)
values ($1, $2)
returning id;

-- name: ProjectGet :one
select *
from unweave.project
where id = $1;

-- name: SessionCreate :one
insert into unweave.session (node_id, created_by, project_id, provider, ssh_key_id)
values ($1, $2, $3, $4, (select id
                         from unweave.ssh_key as ssh_keys
                         where ssh_keys.name = @ssh_key_name
                           and owner_id = $2))
returning id;

-- name: SessionGet :one
select *
from unweave.session
where id = $1;

-- name: SessionAndProjectGet :one
select s.id,
       s.node_id,
       s.created_by,
       s.created_at,
       s.ready_at,
       s.exited_at,
       s.status,
       s.provider,
       s.ssh_key_id,
       p.id         as project_id,
       p.name       as project_name,
       p.owner_id   as project_owner_id,
       p.created_at as project_created_at
from unweave.session s
         join unweave.project p on p.id = s.project_id
where s.id = $1;

-- name: SessionsGet :many
select session.id, ssh_key.name as ssh_key_name, session.status
from unweave.session
         left join unweave.ssh_key
                   on ssh_key.id = session.ssh_key_id
where project_id = $1
order by unweave.session.created_at desc
limit $2 offset $3;

-- name: SessionSetTerminated :exec
update unweave.session
set status = unweave.session_status('terminated')
where id = $1;

-- name: SSHKeyAdd :exec
insert INTO unweave.ssh_key (owner_id, name, public_key)
values ($1, $2, $3);

-- name: SSHKeysGet :many
select *
from unweave.ssh_key
where owner_id = $1;

-- name: SSHKeyGetByName :one
select *
from unweave.ssh_key
where name = $1;

-- name: SSHKeyGetByPublicKey :one
select *
from unweave.ssh_key
where public_key = $1;
