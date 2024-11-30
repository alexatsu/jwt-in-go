-- name: CreateSession :exec
WITH
    deleted_session AS (
        DELETE FROM sessions
        WHERE
            user_id = $1
        RETURNING
            *
    )
INSERT INTO
    sessions (user_id, refresh_token)
VALUES ($1, $2);

-- name: DeleteSession :exec
DELETE FROM sessions WHERE user_id = $1;

-- name: GetSession :one
SELECT * FROM sessions WHERE user_id = $1;