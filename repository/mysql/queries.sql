-- name: GetAll :many
SELECT * FROM players;

-- name: GetOneByUsername :exec
SELECT * FROM players WHERE username = ?;

-- name: CreateUser :exec
INSERT INTO
    players (
        user_uid, username, password, email, created_at
    )
VALUES (?, ?, ?, ?, ?);

-- name: CreateAPIKey :exec
INSERT INTO
    api_keys (api_key, u_id)
SELECT ?, user_uid
FROM players
WHERE
    username = ?;

-- name: GetUserWithAPIKeyByUsername :one
SELECT pl.*, ak.api_key
FROM players pl
    JOIN api_keys ak ON pl.user_uid = ak.u_id
WHERE
    pl.username = ?;

-- name: CreateSystem :exec
INSERT INTO
    systems (
        symbol, sector_symbol, type, x, y
    )
VALUES (?, ?, ?, ?, ?);

-- name: GetSystemBySymbol :exec
SELECT * FROM systems WHERE symbol = ?;