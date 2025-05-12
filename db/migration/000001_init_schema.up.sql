CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE "user" (
  "user_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "username" VARCHAR(50) NOT NULL,
  "email" VARCHAR(100) UNIQUE NOT NULL,
  "password" TEXT NOT NULL,
  "pfp" TEXT NOT NULL DEFAULT '/static/assets/default-avatar-64x64.png',
  "role" VARCHAR(20) NOT NULL DEFAULT 'user',
  "email_verified" BOOL DEFAULT false,
  "banned" BOOL DEFAULT false,
  "is_deleted" BOOL DEFAULT false,
  "created_at" TIMESTAMPTZ DEFAULT (now())
);

CREATE TABLE "session" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID NOT NULL,
  "username" VARCHAR NOT NULL,
  "refresh_token" VARCHAR NOT NULL,
  "user_agent" VARCHAR NOT NULL,
  "client_ip" VARCHAR NOT NULL,
  "is_blocked" BOOLEAN NOT NULL,
  "expires_at" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE "content" (
  "content_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID NOT NULL,
  "category_id" UUID NOT NULL,
  "title" TEXT NOT NULL,
  "thumbnail" TEXT,
  "content_description" TEXT NOT NULL,
  "comments_enabled" BOOL NOT NULL DEFAULT true,
  "view_count_enabled" BOOL NOT NULL DEFAULT true,
  "like_count_enabled" BOOL NOT NULL DEFAULT true,
  "dislike_count_enabled" BOOL NOT NULL DEFAULT false,
  "status" VARCHAR(20) NOT NULL DEFAULT 'draft',
  "view_count" INT NOT NULL DEFAULT 0,
  "like_count" INT NOT NULL DEFAULT 0,
  "dislike_count" INT NOT NULL DEFAULT 0,
  "comment_count" INT NOT NULL DEFAULT 0,
  "created_at" TIMESTAMPTZ DEFAULT (now()),
  "updated_at" TIMESTAMPTZ DEFAULT (now()),
  "published_at" TIMESTAMPTZ DEFAULT NULL,
  "is_deleted" BOOL DEFAULT false
);

CREATE TABLE "views" (
  "view_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "content_id" UUID NOT NULL,
  "user_id" UUID NOT NULL,
  UNIQUE ("content_id", "user_id")
);


CREATE TABLE "category" (
  "category_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "category_name" VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE "tag" (
  "tag_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "tag_name" VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE "content_tag" (
  "content_id" UUID NOT NULL,
  "tag_id" UUID NOT NULL
);

CREATE TABLE "comment" (
  "comment_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "content_id" UUID NOT NULL,
  "user_id" UUID NOT NULL,
  "comment_text" TEXT NOT NULL,
  "score" INT NOT NULL DEFAULT 0,
  "created_at" TIMESTAMPTZ DEFAULT now(),
  "updated_at" TIMESTAMPTZ DEFAULT NULL,
  "is_deleted" BOOL DEFAULT false
);

CREATE TABLE "media" (
  "media_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "content_id" UUID NOT NULL,
  "media_type" VARCHAR(50) NOT NULL,
  "media_url" VARCHAR(255) NOT NULL,
  "media_caption" TEXT NOT NULL,
  "media_order" INT NOT NULL DEFAULT 0
);

CREATE TABLE "content_reaction" (
  "content_id" UUID NOT NULL,
  "user_id" UUID NOT NULL,
  "reaction" VARCHAR(10) NOT NULL,
  CONSTRAINT unique_content_reaction UNIQUE (content_id, user_id)
);

CREATE TABLE "comment_reaction" (
  "comment_id" UUID NOT NULL,
  "user_id" UUID NOT NULL,
  "reaction" VARCHAR(10) NOT NULL,
  CONSTRAINT unique_comment_reaction UNIQUE (comment_id, user_id)
);

CREATE TABLE "ads" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "title" VARCHAR(255),
  "description" TEXT,
  "image_url" VARCHAR(500),
  "target_url" VARCHAR(500),
  "placement" VARCHAR(50),
  "status" VARCHAR(20) CHECK (status IN ('active', 'inactive')),
  "clicks" INT DEFAULT 0,
  "start_date" TIMESTAMPTZ,
  "end_date" TIMESTAMPTZ,
  "created_at" TIMESTAMPTZ DEFAULT (now()),
  "updated_at" TIMESTAMPTZ DEFAULT (now())
);


CREATE TABLE "global_settings" (
  "global_settings_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "disable_comments" BOOL NOT NULL DEFAULT false,
  "disable_likes" BOOL NOT NULL DEFAULT false,
  "disable_dislikes" BOOL NOT NULL DEFAULT false,
  "disable_views" BOOL NOT NULL DEFAULT false,
  "disable_ads" BOOL NOT NULL DEFAULT false
);

CREATE TABLE "analytics_daily" (
  "analytics_date" DATE PRIMARY KEY,
  "total_views" INT NOT NULL DEFAULT 0,
  "total_likes" INT NOT NULL DEFAULT 0,
  "total_dislikes" INT NOT NULL DEFAULT 0,
  "total_comments" INT NOT NULL DEFAULT 0,
  "total_ads_clicks" INT NOT NULL DEFAULT 0,
  "created_at" TIMESTAMPTZ DEFAULT (now()),
  "updated_at" TIMESTAMPTZ DEFAULT (now())
);

-- Create indexes 
CREATE INDEX "idx_user_email" ON "user"("email");
CREATE INDEX "idx_user_username" ON "user"("username");
CREATE INDEX "idx_user_active" ON "user"("is_deleted", "banned", "created_at");
CREATE INDEX "idx_user_banned_created_at" ON "user"("banned", "created_at");
CREATE INDEX "idx_user_deleted_created_at" ON "user"("is_deleted", "created_at");
CREATE INDEX "idx_content_status_published" ON content("status", "is_deleted", "published_at");
CREATE INDEX "idx_content_user" ON content("user_id", "is_deleted");
CREATE INDEX "idx_content_category" ON content("category_id", "status", "is_deleted");
CREATE INDEX "idx_content_fulltext" ON content USING gin (to_tsvector('english', "title" || ' ' || "content_description"));
CREATE INDEX "idx_tag_name_lower" ON tag(lower("tag_name"));
CREATE INDEX "idx_content_tag_content" ON content_tag("content_id");
CREATE INDEX "idx_content_tag_tag" ON content_tag("tag_id");
CREATE INDEX "idx_comment_content_created" ON comment("content_id", "is_deleted", "created_at");
CREATE INDEX "idx_media_content_order" ON media("content_id", "media_order");
CREATE INDEX "idx_content_reaction_content" ON content_reaction("content_id");
CREATE INDEX "idx_comment_reaction_comment" ON comment_reaction("comment_id");
CREATE INDEX "idx_content_reaction_content_user" ON content_reaction("content_id", "user_id");
CREATE INDEX "idx_comment_reaction_comment_user" ON comment_reaction("comment_id", "user_id");
CREATE INDEX "idx_ads_status_start_date_end_date" ON "ads"("status", "start_date", "end_date");
CREATE INDEX "idx_analytics_daily_date" ON "analytics_daily"("analytics_date");
CREATE INDEX "idx_analytics_daily_date_updated_at" ON "analytics_daily"("analytics_date", "updated_at");
CREATE INDEX "idx_recent_analytics" ON "analytics_daily"("analytics_date");

-- Create foreign key constraints
ALTER TABLE "session" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("user_id") ON DELETE CASCADE;
ALTER TABLE "content" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("user_id") ON DELETE CASCADE;
ALTER TABLE "content" ADD FOREIGN KEY ("category_id") REFERENCES "category" ("category_id") ON DELETE CASCADE;
ALTER TABLE "content_tag" ADD FOREIGN KEY ("content_id") REFERENCES "content" ("content_id") ON DELETE CASCADE;
ALTER TABLE "content_tag" ADD FOREIGN KEY ("tag_id") REFERENCES "tag" ("tag_id");
ALTER TABLE "comment" ADD FOREIGN KEY ("content_id") REFERENCES "content" ("content_id") ON DELETE CASCADE;
ALTER TABLE "comment" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("user_id") ON DELETE CASCADE;
ALTER TABLE "media" ADD FOREIGN KEY ("content_id") REFERENCES "content" ("content_id") ON DELETE CASCADE;
ALTER TABLE "content_reaction" ADD FOREIGN KEY ("content_id") REFERENCES "content" ("content_id") ON DELETE CASCADE;
ALTER TABLE "content_reaction" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("user_id") ON DELETE CASCADE;
ALTER TABLE "comment_reaction" ADD FOREIGN KEY ("comment_id") REFERENCES "comment" ("comment_id") ON DELETE CASCADE;
ALTER TABLE "comment_reaction" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("user_id") ON DELETE CASCADE;
ALTER TABLE "comment" ADD COLUMN "parent_comment_id" UUID DEFAULT NULL REFERENCES "comment"("comment_id") ON DELETE CASCADE;

