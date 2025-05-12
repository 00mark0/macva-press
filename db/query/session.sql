-- name: CreateSession :one
INSERT INTO session (
  id,
  user_id,
  username,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, user_id, username, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at;

-- name: GetSession :one
SELECT * FROM session
WHERE id = $1 LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM session
WHERE id = $1;
