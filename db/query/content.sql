-- name: CreateContent :one
INSERT INTO content (
    user_id,
    category_id,
    title,
    content_description,
    comments_enabled,
    view_count_enabled,
    like_count_enabled,
    dislike_count_enabled
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateContent :one
UPDATE content
SET
    title = COALESCE($2, title),
    content_description = COALESCE($3, content_description),
    category_id = COALESCE($4, category_id),
    comments_enabled = COALESCE($5, comments_enabled),
    view_count_enabled = COALESCE($6, view_count_enabled),
    like_count_enabled = COALESCE($7, like_count_enabled),
    dislike_count_enabled = COALESCE($8, dislike_count_enabled),
    updated_at = now()
WHERE content_id = $1
RETURNING *;

-- name: AddThumbnail :one
UPDATE content
SET
    thumbnail = $2,
    updated_at = now()
WHERE content_id = $1
RETURNING *;

-- name: PublishContent :one
UPDATE content
SET
    status = 'published',
    published_at = now(),
    updated_at = now()
WHERE content_id = $1
RETURNING *;

-- name: UnarchiveContent :one
UPDATE content
SET
    status = 'draft',
    is_deleted = false,
    published_at = null,
    updated_at = now()
WHERE content_id = $1
RETURNING *;

-- name: SoftDeleteContent :one
UPDATE content
SET
    is_deleted = true,
    published_at = null,
    updated_at = now()
WHERE content_id = $1
RETURNING *;

-- name: HardDeleteContent :one
DELETE FROM content
WHERE content_id = $1
RETURNING *;

-- name: GetContentDetails :one
SELECT
  c.*,
  u.username,
  cat.category_name,
  (
    SELECT array_agg(t.tag_name)::text[]
    FROM content_tag ct
    JOIN tag t ON ct.tag_id = t.tag_id
    WHERE ct.content_id = c.content_id
  ) AS tags
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.content_id = $1;

-- name: GetPublishedContentCount :one
SELECT count(*)
FROM content
WHERE status = 'published'
  AND is_deleted = false;

-- name: ListPublishedContentLimit :many
SELECT
  c.*,
  u.username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.status = 'published'
  AND c.is_deleted = false
ORDER BY c.published_at DESC
LIMIT $1;

-- name: ListPublishedContentLimitOldest :many
SELECT
  c.*,
  u.username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.status = 'published'
  AND c.is_deleted = false
ORDER BY c.published_at ASC
LIMIT $1;

-- name: ListPublishedContentLimitTitle :many
SELECT
  c.*,
  u.username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.status = 'published'
  AND c.is_deleted = false
ORDER BY c.title ASC
LIMIT $1;

-- name: ListPublishedContent :many
SELECT
  c.*,
  u.username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.status = 'published'
  AND c.is_deleted = false
ORDER BY c.published_at DESC
LIMIT $1 OFFSET $2;

-- name: ListDraftContent :many
SELECT
  c.*,
  u.username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.status = 'draft'
  AND c.is_deleted = false
ORDER BY c.created_at DESC
LIMIT $1;

-- name: ListDraftContentOldest :many
SELECT
  c.*,
  u.username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.status = 'draft'
  AND c.is_deleted = false
ORDER BY c.created_at ASC
LIMIT $1;

-- name: ListDraftContentTitle :many
SELECT
  c.*,
  u.username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.status = 'draft'
  AND c.is_deleted = false
ORDER BY c.title ASC
LIMIT $1;

-- name: ListDeletedContent :many
SELECT
  c.*,
  u.username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.is_deleted = true
ORDER BY c.created_at DESC
LIMIT $1;

-- name: ListDeletedContentOldest :many
SELECT
  c.*,
  u.username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.is_deleted = true
ORDER BY c.created_at ASC
LIMIT $1;

-- name: ListDeletedContentTitle :many
SELECT
  c.*,
  u.username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.is_deleted = true
ORDER BY c.title ASC
LIMIT $1;

-- name: GetContentByCategoryCount :one
SELECT count(*)
FROM content
WHERE category_id = $1
  AND status = 'published'
  AND is_deleted = false;

-- name: ListContentByCategory :many
SELECT
  c.*,
  u.username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.category_id = $1
  AND c.status = 'published'
  AND c.is_deleted = false
ORDER BY c.published_at DESC
LIMIT $2 OFFSET $3;

-- name: ListContentByCategoryLimit :many
SELECT
  c.*,
  u.username
FROM content c
JOIN "user" u ON c.user_id = u.user_id
WHERE c.category_id = $1
  AND c.status = 'published'
  AND c.is_deleted = false
ORDER BY c.published_at DESC
LIMIT $2;

-- name: GetContentByTagCount :one
SELECT count(DISTINCT c.content_id)
FROM content c
JOIN content_tag ct ON c.content_id = ct.content_id
JOIN tag t ON ct.tag_id = t.tag_id
WHERE t.tag_name = $1
  AND c.status = 'published'
  AND c.is_deleted = false;

-- name: ListContentByTag :many
SELECT DISTINCT
  c.*,
  u.username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
JOIN content_tag ct ON c.content_id = ct.content_id
JOIN tag t ON ct.tag_id = t.tag_id
WHERE t.tag_name = $1
  AND c.status = 'published'
  AND c.is_deleted = false
ORDER BY c.published_at DESC
LIMIT $2 OFFSET $3;

-- name: ListContentByTagLimit :many
SELECT DISTINCT
  c.*,
  u.username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
JOIN content_tag ct ON c.content_id = ct.content_id
JOIN tag t ON ct.tag_id = t.tag_id
WHERE t.tag_name = $1
  AND c.status = 'published'
  AND c.is_deleted = false
ORDER BY c.published_at DESC
LIMIT $2;

-- name: GetSearchContentCount :one
SELECT count(DISTINCT c.content_id)
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
LEFT JOIN content_tag ct ON c.content_id = ct.content_id
LEFT JOIN tag t ON ct.tag_id = t.tag_id
WHERE c.status = 'published'
  AND c.is_deleted = false
  AND (
    c.title ILIKE '%' || @search_term::text || '%'
    OR c.content_description ILIKE '%' || @search_term::text || '%'
    OR t.tag_name ILIKE '%' || @search_term::text || '%'
  );

-- name: SearchContent :many
SELECT DISTINCT
  c.*,
  u.username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
LEFT JOIN content_tag ct ON c.content_id = ct.content_id
LEFT JOIN tag t ON ct.tag_id = t.tag_id
WHERE c.status = 'published'
  AND c.is_deleted = false
  AND (
    cat.category_name ILIKE '%' || @search_term::text || '%'
    OR c.title ILIKE '%' || @search_term::text || '%'
    OR c.content_description ILIKE '%' || @search_term::text || '%'
    OR t.tag_name ILIKE '%' || @search_term::text || '%'
  )
ORDER BY c.published_at DESC
LIMIT $1;

-- name: SearchDraftContent :many
SELECT
  c.*,
  u.username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.status = 'draft'
  AND c.is_deleted = false
  AND (
    c.title ILIKE '%' || @search_term::text || '%'
    OR c.content_description ILIKE '%' || @search_term::text || '%'
  )
ORDER BY c.published_at DESC
LIMIT $1;

-- name: SearchDelContent :many
SELECT
  c.*,
  u.username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.is_deleted = true
  AND (
    c.title ILIKE '%' || @search_term::text || '%'
    OR c.content_description ILIKE '%' || @search_term::text || '%'
  )
ORDER BY c.published_at DESC
LIMIT $1;

-- name: IncrementViewCount :one
UPDATE content
SET
  view_count = view_count + 1
WHERE content_id = $1
RETURNING view_count;

-- name: IncrementCommentCount :exec
UPDATE content
SET
  comment_count = comment_count + 1
WHERE content_id = $1;

-- name: InsertOrUpdateContentReaction :one
INSERT INTO content_reaction (content_id, user_id, reaction)
VALUES ($1, $2, $3)
ON CONFLICT (content_id, user_id)
DO UPDATE SET reaction = EXCLUDED.reaction
RETURNING content_id;

-- name: DeleteContentReaction :one
DELETE FROM content_reaction
WHERE content_id = $1 AND user_id = $2
RETURNING content_id;

-- name: UpdateContentLikeDislikeCount :one
UPDATE content c
SET
  like_count = (
    SELECT count(*) 
    FROM content_reaction 
    WHERE content_id = c.content_id AND reaction = 'like'
  ),
  dislike_count = (
    SELECT count(*) 
    FROM content_reaction 
    WHERE content_id = c.content_id AND reaction = 'dislike'
  ),
  updated_at = now()
WHERE c.content_id = $1 
RETURNING *;

-- name: FetchContentReactions :many
SELECT
  cr.*,
  u.username
FROM content_reaction cr
JOIN "user" u ON cr.user_id = u.user_id
WHERE cr.content_id = $1
LIMIT $2;

-- name: FetchUserContentReaction :one
SELECT
  cr.*,
  u.username
FROM content_reaction cr
JOIN "user" u ON cr.user_id = u.user_id
WHERE cr.content_id = $1 AND cr.user_id = $2
LIMIT 1;

-- name: ListTrendingContent :many
SELECT 
  c.*,
  cat.category_name,
  (c.view_count + c.like_count + c.comment_count) AS total_interactions
FROM content c
JOIN category cat ON c.category_id = cat.category_id
WHERE c.status = 'published'
  AND c.is_deleted = false
  AND c.published_at >= $1
ORDER BY total_interactions DESC
LIMIT $2;

-- name: GetContentOverview :one
SELECT 
  COUNT(*) FILTER (WHERE status = 'draft' AND is_deleted = false) AS draft_count,
  COUNT(*) FILTER (WHERE status = 'published' AND is_deleted = false) AS published_count,
  COUNT(*) FILTER (WHERE is_deleted = true) AS deleted_count
FROM content;

-- name: ListRelatedContent :many
SELECT c.*
FROM content c
WHERE c.content_id <> $1
  AND c.status = 'published'
  AND c.is_deleted = false
  AND c.category_id = (SELECT category_id FROM content WHERE content_id = $1)
  AND EXISTS (
      SELECT 1
      FROM content_tag ct
      WHERE ct.content_id = c.content_id
        AND ct.tag_id IN (
            SELECT tag_id
            FROM content_tag
            WHERE content_id = $1
        )
  )
ORDER BY c.published_at DESC
LIMIT $2;














