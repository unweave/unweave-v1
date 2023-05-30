-- name: ExecListByProvider :many
select *
from unweave.exec as e
where e.provider = $1;

-- name: ExecListActiveByProvider :many
select *
from unweave.exec as e
where provider = $1
  and (status = 'initializing'
    or status = 'running');
