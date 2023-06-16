-- name: VolumeCreate :one
insert into unweave.volume (id, project_id, provider, name, size)
values($1, $2, $3, $4, $5)
returning *;

-- name: VolumeDelete :exec
delete from unweave.volume
where id = $1;

-- name: VolumeGet :one
select * from unweave.volume
where project_id = $1 and (id = $2 or name = $2);

-- name: VolumeList :many
select * from unweave.volume
where project_id = $1;

-- name: VolumeUpdate :exec
update unweave.volume
set size = $2
where id = $1;
