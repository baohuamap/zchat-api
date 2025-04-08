
-- Create "friendship" table
CREATE TABLE "public"."friendship" (
    "id" bigserial NOT NULL,
    "created_at" timestamptz NULL,
    "updated_at" timestamptz NULL,
    "deleted_at" timestamptz NULL,
    "user_id" bigint NULL,
    "friend_id" bigint NULL,
    "status" text NULL,
    PRIMARY KEY ("id")
);

-- Create index "idx_friendship_deleted_at" to table: "friendship"
CREATE INDEX "idx_friendship_deleted_at" ON "public"."friendship" ("deleted_at");

-- Create index "idx_friendship_user_id" to table: "friendship"
CREATE INDEX "idx_friendship_user_id" ON "public"."friendship" ("user_id");

-- Create index "idx_friendship_friend_id" to table: "friendship"
CREATE INDEX "idx_friendship_friend_id" ON "public"."friendship" ("friend_id");



---- create above / drop below ----

DROP TABLE friendship CASCADE;




