-- name: GetAll :many
SELECT * FROM users;

-- name: GetOneByUsername :one
SELECT * FROM users WHERE username = $1;