-- name: GetAll :many
SELECT * FROM players;

-- name: GetOneByUsername :one
SELECT * FROM players WHERE username = $1;

-- name: CreateUser :one
INSERT INTO
    players (
        user_uid, username, password, email, created_at
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING
    user_uid;

-- name: CreateAPIKey :exec
INSERT INTO
    api_keys (key, u_id)
SELECT $1, user_uid
FROM players
WHERE
    username = $2;

-- name: GetUserWithAPIKeyByUsername :one
SELECT pl.*, ak.key
FROM players pl
    JOIN api_keys ak ON pl.user_uid = ak.u_id
WHERE
    pl.username = $1;