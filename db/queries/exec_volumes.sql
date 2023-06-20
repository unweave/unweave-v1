-- name: ExecVolumeGet :many
select exec_id, volume_id, mount_path from unweave.exec_volume as ev where ev.exec_id = $1;

-- name: ExecVolumeCreate :exec
insert into unweave.exec_volume (exec_id, volume_id, mount_path)
values ($1, $2, $3);

-- name: ExecVolumeDelete :exec
delete from unweave.exec_volume
where exec_id = $1;
