-- name: CreateGlobalSettings :one
INSERT INTO "global_settings" ("disable_comments", "disable_likes", "disable_dislikes", "disable_views", "disable_ads")
VALUES (false, false, true, false, false)
RETURNING *;

-- name: GetGlobalSettings :many
SELECT * FROM "global_settings";

-- name: UpdateGlobalSettings :exec
UPDATE "global_settings"
SET
    "disable_comments" = $1,
    "disable_likes" = $2,
    "disable_dislikes" = $3,
    "disable_views" = $4,
    "disable_ads" = $5
WHERE "global_settings_id" = (SELECT "global_settings_id" FROM "global_settings" LIMIT 1);

-- name: ResetGlobalSettings :exec
UPDATE "global_settings"
SET
    "disable_comments" = false,
    "disable_likes" = false,
    "disable_dislikes" = true,
    "disable_views" = false,
    "disable_ads" = false
WHERE "global_settings_id" = (SELECT "global_settings_id" FROM "global_settings" LIMIT 1);

