-- name: ExecSSHKeyInsert :exec
INSERT INTO unweave.exec_ssh_key (exec_id, ssh_key_id)
VALUES ($1, $2);

-- name: ExecSSHKeyGet :one
SELECT *
FROM unweave.exec_ssh_key
WHERE exec_id = $1
  AND ssh_key_id = $2;

-- name: ExecSSHKeysByExecIDGet :many
SELECT *
FROM unweave.exec_ssh_key
WHERE exec_id = $1;

-- name: ExecSSHKeyDelete :exec
DELETE FROM unweave.exec_ssh_key
WHERE exec_id = $1
  AND ssh_key_id = $2;
