-- name: CreateTag :one
INSERT INTO tag (tag_name)
VALUES ($1)
RETURNING tag_id, tag_name;

-- name: UpdateTag :one
UPDATE tag
SET tag_name = $1
WHERE tag_id = $2
RETURNING tag_id, tag_name;

-- name: DeleteTag :exec
WITH deleted_tag AS (
  DELETE FROM tag
  WHERE tag.tag_id = $1
  RETURNING tag_id
)
DELETE FROM content_tag
WHERE tag_id IN (SELECT tag_id FROM deleted_tag);

-- name: GetTag :one
SELECT tag_id, tag_name
FROM tag
WHERE tag_id = $1;

-- name: ListTags :many
SELECT tag_id, tag_name
FROM tag
ORDER BY tag_name ASC
LIMIT $1;

-- name: SearchTags :many
SELECT tag_id, tag_name
FROM tag
WHERE lower(tag_name) LIKE lower(@search::text)
ORDER BY tag_name ASC
LIMIT $1;

-- name: AddTagToContent :exec
INSERT INTO content_tag (content_id, tag_id)
VALUES ($1, $2);

-- name: GetTagsByContent :many
SELECT tag.tag_id, tag.tag_name
FROM tag
JOIN content_tag ct ON tag.tag_id = ct.tag_id
WHERE ct.content_id = $1;

-- name: GetUniqueTagsByCategoryID :many
SELECT DISTINCT t.tag_id, t.tag_name
FROM tag t
JOIN content_tag ct ON t.tag_id = ct.tag_id
JOIN content c ON ct.content_id = c.content_id
WHERE c.category_id = $1
  AND c.status = 'published'
  AND c.is_deleted = false
ORDER BY t.tag_name;

-- name: RemoveTagFromContent :exec
DELETE FROM content_tag
WHERE content_id = $1 AND tag_id = $2;


