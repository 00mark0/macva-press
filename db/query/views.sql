-- name: GetView :one
SELECT 1
FROM "views"
WHERE "content_id" = $1
  AND "user_id" = $2
LIMIT 1;

-- name: AddView :exec
INSERT INTO "views" ("content_id", "user_id")
VALUES ($1, $2);
