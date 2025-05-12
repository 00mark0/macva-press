-- Drop foreign key constraints
ALTER TABLE "content" DROP CONSTRAINT IF EXISTS "content_user_id_fkey";
ALTER TABLE "content" DROP CONSTRAINT IF EXISTS "content_category_id_fkey";
ALTER TABLE "content_tag" DROP CONSTRAINT IF EXISTS "content_tag_content_id_fkey";
ALTER TABLE "content_tag" DROP CONSTRAINT IF EXISTS "content_tag_tag_id_fkey";
ALTER TABLE "comment" DROP CONSTRAINT IF EXISTS "comment_content_id_fkey";
ALTER TABLE "comment" DROP CONSTRAINT IF EXISTS "comment_user_id_fkey";
ALTER TABLE "media" DROP CONSTRAINT IF EXISTS "media_content_id_fkey";
ALTER TABLE "content_reaction" DROP CONSTRAINT IF EXISTS "content_reaction_content_id_fkey";
ALTER TABLE "content_reaction" DROP CONSTRAINT IF EXISTS "content_reaction_user_id_fkey";
ALTER TABLE "comment_reaction" DROP CONSTRAINT IF EXISTS "comment_reaction_comment_id_fkey";
ALTER TABLE "comment_reaction" DROP CONSTRAINT IF EXISTS "comment_reaction_user_id_fkey";
ALTER TABLE "user" DROP CONSTRAINT IF EXISTS "user_role_id_fkey";

-- Drop indexes
DROP INDEX IF EXISTS "idx_user_email";
DROP INDEX IF EXISTS "idx_user_username";
DROP INDEX IF EXISTS "idx_content_user_id";
DROP INDEX IF EXISTS "idx_content_category_id";
DROP INDEX IF EXISTS "idx_content_created_at";
DROP INDEX IF EXISTS "idx_content_updated_at";
DROP INDEX IF EXISTS "idx_category_name";
DROP INDEX IF EXISTS "idx_tag_name";
DROP INDEX IF EXISTS "idx_content_tag_tag_id_content_id";
DROP INDEX IF EXISTS "idx_comment_content_id";
DROP INDEX IF EXISTS "idx_comment_user_id";
DROP INDEX IF EXISTS "idx_comment_created_at";
DROP INDEX IF EXISTS "idx_media_content_id";
DROP INDEX IF EXISTS "idx_media_order";
DROP INDEX IF EXISTS "idx_content_reaction_reaction";
DROP INDEX IF EXISTS "idx_comment_reaction_reaction";
DROP INDEX IF EXISTS "idx_analytics_daily_created_at";

-- Drop tables
DROP TABLE IF EXISTS "content_reaction" CASCADE;
DROP TABLE IF EXISTS "comment_reaction" CASCADE;
DROP TABLE IF EXISTS "media" CASCADE;
DROP TABLE IF EXISTS "comment" CASCADE;
DROP TABLE IF EXISTS "content_tag" CASCADE;
DROP TABLE IF EXISTS "tag" CASCADE;
DROP TABLE IF EXISTS "category" CASCADE;
DROP TABLE IF EXISTS "content" CASCADE;
DROP TABLE IF EXISTS "role" CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;
DROP TABLE IF EXISTS "global_settings" CASCADE;
DROP TABLE IF EXISTS "analytics_daily" CASCADE;
DROP TABLE IF EXISTS "ads" CASCADE;
DROP TABLE IF EXISTS "views" CASCADE;
DROP TABLE IF EXISTS "session" CASCADE;

