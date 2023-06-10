-- name: ProjectGet :one
select *
from unweave.project
where id = $1;
