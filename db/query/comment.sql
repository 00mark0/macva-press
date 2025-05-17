-- name: CreateComment :one
INSERT INTO comment (content_id, user_id, comment_text)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateComment :one
UPDATE comment
SET 
  comment_text = $1,
  updated_at = now()
WHERE 
  comment_id = $2
  AND is_deleted = false
RETURNING
  comment_id,
  comment_text,
  updated_at;

-- name: SoftDeleteComment :one
UPDATE comment
SET 
  is_deleted = true,
  updated_at = now()
WHERE 
  comment_id = $1
  AND is_deleted = false
RETURNING
  comment_id,
  is_deleted,
  updated_at;

-- name: DeleteComment :one
DELETE FROM comment
WHERE comment_id = $1
RETURNING *;

-- name: ListContentComments :many
SELECT
  cm.*,
  u.username,
  u.pfp,
  u.role
FROM comment cm
JOIN "user" u ON cm.user_id = u.user_id
WHERE cm.content_id = $1
  AND cm.is_deleted = false
  AND cm.parent_comment_id IS NULL
ORDER BY cm.created_at DESC
LIMIT $2;

-- name: ListContentCommentsByScore :many
SELECT
  cm.*,
  u.username,
  u.pfp,
  u.role
FROM comment cm
JOIN "user" u ON cm.user_id = u.user_id
WHERE cm.content_id = $1
  AND cm.is_deleted = false
  AND cm.parent_comment_id IS NULL
ORDER BY cm.score DESC
LIMIT $2;

-- name: InsertOrUpdateCommentReaction :one
INSERT INTO comment_reaction (comment_id, user_id, reaction)
VALUES ($1, $2, $3)
ON CONFLICT (comment_id, user_id)
DO UPDATE SET reaction = EXCLUDED.reaction
RETURNING comment_id;

-- name: DeleteCommentReaction :one
DELETE FROM comment_reaction
WHERE comment_id = $1 AND user_id = $2
RETURNING comment_id;

-- name: UpdateCommentScore :one
UPDATE comment c
SET
  score = (
    SELECT count(*) 
    FROM comment_reaction 
    WHERE comment_id = c.comment_id AND reaction = 'like'
  )
  -
  (
    SELECT count(*) 
    FROM comment_reaction 
    WHERE comment_id = c.comment_id AND reaction = 'dislike'
  ),
  updated_at = c.updated_at  -- Explicitly keep the current value
WHERE c.comment_id = $1
RETURNING comment_id, content_id, user_id, comment_text, score, created_at, updated_at, is_deleted, parent_comment_id;

-- name: FetchCommentReactions :many
SELECT
  cr.*,
  u.username
FROM comment_reaction cr
JOIN "user" u ON cr.user_id = u.user_id
WHERE cr.comment_id = $1
LIMIT $2;

-- name: GetUserCommentReaction :one
SELECT * FROM comment_reaction
WHERE comment_id = $1 AND user_id = $2;

-- name: GetCommentByID :one
SELECT * FROM comment
WHERE comment_id = $1;

-- name: GetUserReactionsForContentComments :many
SELECT
  cr.*,
  u.username
FROM comment_reaction cr
JOIN comment c ON cr.comment_id = c.comment_id
JOIN "user" u ON cr.user_id = u.user_id
WHERE c.content_id = $1
  AND cr.user_id = $2
  AND c.is_deleted = false;

-- name: GetCommentCountForContent :one
SELECT count(*) FROM comment
WHERE content_id = $1
  AND is_deleted = false;

-- name: CreateReply :one
INSERT INTO comment (content_id, user_id, comment_text, parent_comment_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListCommentReplies :many
SELECT cm.*, u.username, u.pfp, u.role 
FROM comment cm 
JOIN "user" u ON cm.user_id = u.user_id
WHERE cm.parent_comment_id = $1 AND cm.is_deleted = false
ORDER BY cm.created_at ASC
LIMIT $2;

-- name: GetReplyCount :one
SELECT COUNT(*) 
FROM comment
WHERE parent_comment_id = $1 AND is_deleted = false;
