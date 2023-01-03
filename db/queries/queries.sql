-- name: ProjectCreate :exec
insert into unweave.projects (name, owner_id)
values ($1, $2);

-- name: ProjectGet :one
select *
from unweave.projects
where id = $1;

-- name: SessionCreate :exec
insert into unweave.sessions (node_id, created_by, project_id, runtime)
values ($1, $2, $3, $4);

-- name: SessionGet :one
select *
from unweave.sessions
where id = $1;

-- name: SessionSetTerminated :exec
update unweave.sessions
set status = unweave.session_status('terminated')
where id = $1;

-- name: SSHKeyAdd :exec
insert INTO unweave.ssh_keys (owner_id, name, public_key)
values ($1, $2, $3);

-- name: SSHKeyGetByName :one
select *
from unweave.ssh_keys
where name = $1;

-- name: SSHKeyGetByPublicKey :one
select *
from unweave.ssh_keys
where public_key = $1;
