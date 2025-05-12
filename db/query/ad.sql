-- name: CreateAd :one
INSERT INTO "ads"
("title", "description", "image_url", "target_url", "placement", "status", "start_date", "end_date")
VALUES
($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateAd :one
UPDATE "ads"
SET
  "title" = $1,
  "description" = $2,
  "image_url" = $3,
  "target_url" = $4,
  "placement" = $5,
  "status" = $6,
  "start_date" = $7,
  "end_date" = $8,
  "updated_at" = now()
WHERE "id" = $9
RETURNING *;

-- name: DeactivateAd :one
UPDATE "ads"
SET
  "status" = 'inactive', 
  "updated_at" = now(),
  "start_date" = NULL,
  "end_date" = NULL
WHERE "id" = $1
RETURNING *;

-- name: DeleteAd :exec
DELETE FROM "ads"
WHERE "id" = $1;

-- name: GetAd :one
SELECT *
FROM "ads"
WHERE "id" = $1;

-- name: ListAds :many
SELECT *
FROM "ads"
LIMIT $1;

-- name: ListInactiveAds :many
SELECT *
FROM "ads"
WHERE "status" = 'inactive'
ORDER BY "created_at" DESC
LIMIT $1;

-- name: ListScheduledAds :many
SELECT *
FROM "ads"
WHERE "status" = 'active'
  AND "start_date" > now() AT TIME ZONE 'Europe/Belgrade'
ORDER BY "start_date" ASC
LIMIT $1;

-- name: ListActiveAds :many
SELECT *
FROM "ads"
WHERE "status" = 'active'
  AND "start_date" <= now() AT TIME ZONE 'Europe/Belgrade'
  AND "end_date" >= now() AT TIME ZONE 'Europe/Belgrade'
ORDER BY "created_at" DESC
LIMIT $1;

-- name: ListAdsByPlacement :many
SELECT *
FROM "ads"
WHERE "placement" = $1
  AND "status" = 'active'
  AND "start_date" <= now() AT TIME ZONE 'Europe/Belgrade' 
  AND "end_date" >= now() AT TIME ZONE 'Europe/Belgrade'
LIMIT $2;

-- name: IncrementAdClicks :one
UPDATE "ads"
SET
  "clicks" = "clicks" + 1, 
  "updated_at" = now()
WHERE "id" = $1
RETURNING "clicks";
