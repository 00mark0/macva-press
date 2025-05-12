-- name: CreateUser :one
INSERT INTO "user" (username, email, password) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateUserAdmin :one
INSERT INTO "user" (username, email, password, role) 
VALUES ($1, $2, $3, $4)  
RETURNING *;

-- name: GetUserByID :one
SELECT user_id, username, password, email, role, pfp, email_verified, banned, is_deleted 
FROM "user" 
WHERE user_id = $1 AND banned = false;

-- name: GetUserByEmail :one
SELECT user_id, username, password, email, role, pfp, email_verified, banned, is_deleted 
FROM "user" 
WHERE email = $1 AND banned = false;

-- name: UpdateUser :exec
UPDATE "user" 
SET username = COALESCE($2, username),
    pfp = COALESCE($3, pfp)
WHERE user_id = $1;

-- name: UpdateUserPassword :exec
UPDATE "user" 
SET password = $2 
WHERE user_id = $1;

-- name: BanUser :exec
UPDATE "user" 
SET banned = true 
WHERE user_id = $1;

-- name: UnbanUser :exec
UPDATE "user" 
SET banned = false 
WHERE user_id = $1;

-- name: DeleteUser :exec
UPDATE "user" 
SET email = CONCAT('deleted_', user_id, '@example.com'), 
    password = '', 
    pfp = '/static/assets/default-avatar-64x64.png', 
    is_deleted = true 
WHERE user_id = $1;

-- name: CheckEmailExists :one
SELECT 1 
FROM "user" 
WHERE email = $1;

-- name: GetAdminUsers :many
SELECT user_id, username, email, password, pfp, role, email_verified, banned, is_deleted, created_at
FROM "user"
WHERE "role" = 'admin'
ORDER BY created_at DESC;

-- name: GetActiveUsersCount :one
SELECT COUNT(*) AS count
FROM "user"
WHERE "is_deleted" = false
  AND "banned" = false;

-- name: GetActiveUsers :many
SELECT user_id, username, email, password, pfp, role, email_verified, banned, is_deleted, created_at
FROM "user"
WHERE "is_deleted" = false
  AND "banned" = false
ORDER BY created_at DESC
LIMIT $1;

-- name: GetActiveUsersOldest :many
SELECT user_id, username, email, password, pfp, role, email_verified, banned, is_deleted, created_at
FROM "user"
WHERE "is_deleted" = false
  AND "banned" = false
ORDER BY created_at ASC
LIMIT $1;

-- name: GetActiveUsersTitle :many
SELECT user_id, username, email, password, pfp, role, email_verified, banned, is_deleted, created_at
FROM "user"
WHERE "is_deleted" = false
  AND "banned" = false
ORDER BY username ASC
LIMIT $1;

-- name: GetBannedUsersCount :one
SELECT COUNT(*) AS count
FROM "user"
WHERE "banned" = true AND "is_deleted" = false;

-- name: GetBannedUsers :many
SELECT user_id, username, email, password, pfp, role, email_verified, banned, is_deleted, created_at
FROM "user"
WHERE "banned" = true
  AND "is_deleted" = false
ORDER BY created_at DESC
LIMIT $1;

-- name: GetBannedUsersOldest :many
SELECT user_id, username, email, password, pfp, role, email_verified, banned, is_deleted, created_at
FROM "user"
WHERE "banned" = true
  AND "is_deleted" = false
ORDER BY created_at ASC
LIMIT $1;

-- name: GetBannedUsersTitle :many
SELECT user_id, username, email, password, pfp, role, email_verified, banned, is_deleted, created_at
FROM "user"
WHERE "banned" = true
  AND "is_deleted" = false
ORDER BY username ASC
LIMIT $1;

-- name: GetDeletedUsersCount :one
SELECT COUNT(*) AS count
FROM "user"
WHERE "is_deleted" = true;

-- name: GetDeletedUsers :many
SELECT user_id, username, email, password, pfp, role, email_verified, banned, is_deleted, created_at
FROM "user"
WHERE "is_deleted" = true
ORDER BY created_at DESC
LIMIT $1;

-- name: GetDeletedUsersOldest :many
SELECT user_id, username, email, password, pfp, role, email_verified, banned, is_deleted, created_at
FROM "user"
WHERE "is_deleted" = true
ORDER BY created_at ASC
LIMIT $1;

-- name: GetDeletedUsersTitle :many
SELECT user_id, username, email, password, pfp, role, email_verified, banned, is_deleted, created_at
FROM "user"
WHERE "is_deleted" = true
ORDER BY username ASC
LIMIT $1;

-- name: SearchActiveUsers :many
SELECT
  u.*
FROM "user" u
WHERE u.is_deleted = false 
  AND u.banned = false
  AND (
    u.username ILIKE '%' || @search_term::text || '%'
    OR u.email ILIKE '%' || @search_term::text || '%'
  )
ORDER BY u.created_at DESC
LIMIT $1;

-- name: SearchDeletedUsers :many
SELECT
  u.*
FROM "user" u
WHERE u.is_deleted = true
  AND (
    u.username ILIKE '%' || @search_term::text || '%'
    OR u.email ILIKE '%' || @search_term::text || '%'
  )
ORDER BY u.created_at DESC
LIMIT $1;

-- name: SearchBannedUsers :many
SELECT
  u.*
FROM "user" u
WHERE u.banned = true
  AND (
    u.username ILIKE '%' || @search_term::text || '%'
    OR u.email ILIKE '%' || @search_term::text || '%'
  )
ORDER BY u.created_at DESC
LIMIT $1;

-- name: SetEmailVerified :exec
UPDATE "user" 
SET email_verified = true 
WHERE user_id = $1;

-- name: CheckAdminExists :one
SELECT EXISTS (
  SELECT 1
  FROM "user"
  WHERE "role" = 'admin'
);








