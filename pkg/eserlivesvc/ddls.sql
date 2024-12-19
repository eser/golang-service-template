CREATE TABLE IF NOT EXISTS "user" (
    "id" CHAR(26) PRIMARY KEY,
    "kind" TEXT NOT NULL CHECK ("kind" IN ('admin', 'editor', 'regular')),
    "name" TEXT NOT NULL,
    "email" TEXT UNIQUE,
    "phone" TEXT,
    "individual_profile_id" CHAR(26),
    "github_remote_id" TEXT UNIQUE,
    "github_handle" TEXT,
    "x_remote_id" TEXT,
    "x_handle" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "profile" (
    "id" CHAR(26) PRIMARY KEY,
    "kind" TEXT NOT NULL CHECK ("kind" IN ('individual', 'organization', 'product', 'venue')),
    "slug" TEXT NOT NULL UNIQUE,
    "profile_picture_uri" TEXT,
    "title" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "show_stories" BOOLEAN NOT NULL DEFAULT FALSE,
    "show_projects" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "profile_membership" (
    "id" CHAR(26) PRIMARY KEY,
    "kind" TEXT NOT NULL CHECK ("kind" IN ('owner', 'lead', 'maintainer', 'contributor', 'sponsor', 'follower')),
    "profile_id" CHAR(26) NOT NULL,
    "user_id" CHAR(26) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ,
    UNIQUE("profile_id", "user_id")
);

CREATE TABLE IF NOT EXISTS "event_series" (
    "id" CHAR(26) PRIMARY KEY,
    "slug" TEXT NOT NULL UNIQUE,
    "event_picture_uri" TEXT,
    "title" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "event" (
    "id" CHAR(26) PRIMARY KEY,
    "kind" TEXT NOT NULL CHECK ("kind" IN ('broadcast', 'meetup', 'conference')),
    "status" TEXT NOT NULL CHECK ("status" IN ('draft', 'published', 'archived')),
    "series_id" CHAR(26),
    "slug" TEXT NOT NULL UNIQUE,
    "event_picture_uri" TEXT,
    "title" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "attendance_uri" TEXT,
    "published_at" TIMESTAMPTZ,
    "time_start" TIMESTAMPTZ NOT NULL,
    "time_end" TIMESTAMPTZ NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "event_attendance" (
    "id" CHAR(26) PRIMARY KEY,
    "kind" TEXT NOT NULL CHECK ("kind" IN ('organizer', 'co-organizer', 'speaker', 'sponsor', 'guest')),
    "event_id" CHAR(26) NOT NULL,
    "profile_id" CHAR(26) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ,
    UNIQUE("event_id", "profile_id")
);

CREATE TABLE IF NOT EXISTS "question" (
    "id" CHAR(26) PRIMARY KEY,
    "user_id" CHAR(26) NOT NULL,
    "content" TEXT NOT NULL,
    "answer_kind" TEXT CHECK ("answer_kind" IN ('text', 'article', 'video')),
    "answer_content" TEXT,
    "answer_uri" TEXT,
    "answered_at" TIMESTAMPTZ,
    "is_anonymous" BOOLEAN NOT NULL DEFAULT FALSE,
    "is_hidden" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "question_vote" (
    "id" CHAR(26) PRIMARY KEY,
    "question_id" CHAR(26) NOT NULL,
    "user_id" CHAR(26) NOT NULL,
    "score" INTEGER NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE("question_id", "user_id")
);

CREATE TABLE IF NOT EXISTS "session" (
    "id" CHAR(26) PRIMARY KEY,
    "status" TEXT NOT NULL CHECK ("status" IN ('active', 'logged_out', 'expired', 'revoked', 'login_requested')),
    "oauth_request_state" TEXT NOT NULL,
    "oauth_request_code_verifier" TEXT NOT NULL,
    "oauth_redirect_uri" TEXT,
    "logged_in_user_id" CHAR(26),
    "logged_in_at" TIMESTAMPTZ,
    "expires_at" TIMESTAMPTZ,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "story" (
    "id" CHAR(26) PRIMARY KEY,
    "kind" TEXT NOT NULL CHECK ("kind" IN ('status', 'announcement', 'news', 'article', 'recipe')),
    "status" TEXT NOT NULL CHECK ("status" IN ('draft', 'published', 'archived')),
    "is_featured" BOOLEAN DEFAULT FALSE,
    "slug" TEXT NOT NULL UNIQUE,
    "story_picture_uri" TEXT,
    "title" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "author_profile_id" CHAR(26),
    "summary" TEXT NOT NULL,
    "content" TEXT NOT NULL,
    "published_at" TIMESTAMPTZ,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ,
    "deleted_at" TIMESTAMPTZ
);

-- Add foreign key constraints
ALTER TABLE "user" ADD FOREIGN KEY ("individual_profile_id") REFERENCES "profile"("id");
ALTER TABLE "profile_membership" ADD FOREIGN KEY ("profile_id") REFERENCES "profile"("id");
ALTER TABLE "profile_membership" ADD FOREIGN KEY ("user_id") REFERENCES "user"("id");
ALTER TABLE "event" ADD FOREIGN KEY ("series_id") REFERENCES "event_series"("id");
ALTER TABLE "event_attendance" ADD FOREIGN KEY ("event_id") REFERENCES "event"("id");
ALTER TABLE "event_attendance" ADD FOREIGN KEY ("profile_id") REFERENCES "profile"("id");
ALTER TABLE "question" ADD FOREIGN KEY ("user_id") REFERENCES "user"("id");
ALTER TABLE "question_vote" ADD FOREIGN KEY ("question_id") REFERENCES "question"("id");
ALTER TABLE "question_vote" ADD FOREIGN KEY ("user_id") REFERENCES "user"("id");
ALTER TABLE "session" ADD FOREIGN KEY ("logged_in_user_id") REFERENCES "user"("id");
ALTER TABLE "story" ADD FOREIGN KEY ("author_profile_id") REFERENCES "profile"("id");

-- Add indexes
CREATE INDEX "session_logged_in_user_id_idx" ON "session"("logged_in_user_id");
