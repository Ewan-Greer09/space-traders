-- name: GetAll :many
SELECT * FROM users;

-- name: GetOneByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: CreateUser :one
INSERT INTO
    users (username, password, api_key)
VALUES ($1, $2, $3)
RETURNING
    id,
    username,
    password,
    api_key;