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
