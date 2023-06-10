-- name: VolumeCreate :one
insert into unweave.volume (id, project_id, provider)
values($1, $2, $3)
returning *;

-- name: VolumeDelete :exec
delete from unweave.volume
where id = $1;

-- name: VolumeGet :one
select * from unweave.volume
where id = $1;

-- name: VolumeList :many
select * from unweave.volume
where project_id = $1;