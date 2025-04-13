
-- Create "friendships" table
CREATE TABLE "public"."friendships" (
    "id" bigserial NOT NULL,
    "created_at" timestamptz NULL,
    "updated_at" timestamptz NULL,
    "deleted_at" timestamptz NULL,
    "user_id" bigint NULL,
    "friend_id" bigint NULL,
    "status" text NULL,
    PRIMARY KEY ("id")
);

-- Create index "idx_friendship_deleted_at" to table: "friendships"
CREATE INDEX "idx_friendship_deleted_at" ON "public"."friendships" ("deleted_at");

-- Create index "idx_friendship_user_id" to table: "friendships"
CREATE INDEX "idx_friendship_user_id" ON "public"."friendships" ("user_id");

-- Create index "idx_friendship_friend_id" to table: "friendships"
CREATE INDEX "idx_friendship_friend_id" ON "public"."friendships" ("friend_id");



---- create above / drop below ----

DROP TABLE friendships CASCADE;




