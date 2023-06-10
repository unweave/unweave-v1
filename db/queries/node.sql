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
