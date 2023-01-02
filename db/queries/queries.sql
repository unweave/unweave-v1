-- name: SessionGet :one
SELECT * FROM unweave.sessions WHERE id = $1;


-- name: SSHKeyAdd :exec
INSERT INTO unweave.ssh_keys (owner_id, name, public_key) VALUES ($1, $2, $3);

-- name: SSHKeyGetByName :one
SELECT * FROM unweave.ssh_keys WHERE name = $1;

-- name: SSHKeyGetByPublicKey :one
SELECT * FROM unweave.ssh_keys WHERE public_key = $1;


