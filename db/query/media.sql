-- name: InsertMedia :one
INSERT INTO media (content_id, media_type, media_url, media_caption, media_order)
VALUES ($1, $2, $3, $4, $5)
RETURNING media_id, content_id, media_type, media_url, media_caption, media_order;

-- name: UpdateMedia :one
UPDATE media
SET media_url = $1,
    media_caption = $2,
    media_order = $3
WHERE media_id = $4
RETURNING media_id, content_id, media_type, media_url, media_caption, media_order;

-- name: DeleteMedia :exec
DELETE FROM media
WHERE media_id = $1;

-- name: GetMediaByID :one
SELECT media_id, content_id, media_type, media_url, media_caption, media_order
FROM media
WHERE media_id = $1;

-- name: ListMediaForContent :many
SELECT media_id, content_id, media_type, media_url, media_caption, media_order
FROM media
WHERE content_id = $1
ORDER BY media_order ASC;

-- name: BatchUpdateMediaOrder :exec
UPDATE media
SET media_order = data.new_order
FROM (
    VALUES 
      -- Format: (media_id, new_order)
      (@media1_id::uuid, @media1_order::int),
      (@media2_id::uuid, @media2_order::int)
      -- Add more tuples as needed...
) AS data(media_id, new_order)
WHERE media.media_id = data.media_id;


