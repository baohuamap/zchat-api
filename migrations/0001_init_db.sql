
-- Create "users" table
CREATE TABLE "public"."users" (
    "id" bigserial NOT NULL,
    "created_at" timestamptz NULL,
    "updated_at" timestamptz NULL,
    "deleted_at" timestamptz NULL,
    "username" text NOT NULL,
    "password" text NOT NULL,
    "email" text UNIQUE NOT NULL,
    "first_name" text NULL,
    "last_name" text NULL,
    "avatar" text NULL,
    "phone" text UNIQUE NOT NULL,
    PRIMARY KEY ("id")
);

-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");

-- Create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "public"."users" ("username");

-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");


-- Create "conversation_type" enum type
CREATE TYPE "conversation_type" AS ENUM ('private', 'group');

-- Create "conversations" table
CREATE TABLE "public"."conversations" (
    "id" bigserial NOT NULL,
    "created_at" timestamptz NULL,
    "updated_at" timestamptz NULL,
    "deleted_at" timestamptz NULL,
    "type" "conversation_type" NOT NULL,
    "creator_id" bigint NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_conversations_creator_id" FOREIGN KEY ("creator_id") REFERENCES "users"("id")
);

-- Create index "idx_files_deleted_at" to table: "files"
CREATE INDEX "idx_conversation_deleted_at" ON "public"."conversations" ("deleted_at");

-- Create "participants" table
CREATE TABLE "public"."participants" (
    "id" bigserial NOT NULL,
    "created_at" timestamptz NULL,
    "updated_at" timestamptz NULL,
    "deleted_at" timestamptz NULL,
    "user_id" bigint NULL,
    "conversation_id" bigint NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_participants_user_id" FOREIGN KEY ("user_id") REFERENCES "users"("id"),
    CONSTRAINT "fk_participants_conversation_id" FOREIGN KEY ("conversation_id") REFERENCES "conversations"("id")
);
-- Create index "idx_participants_deleted_at" to table: "participants"
CREATE INDEX "idx_participants_deleted_at" ON "public"."participants" ("deleted_at");


-- Create "messages" table
CREATE TABLE "public"."messages" (
    "id" bigserial NOT NULL,
    "created_at" timestamptz NULL,
    "updated_at" timestamptz NULL,
    "deleted_at" timestamptz NULL,
    "sender_id" bigint NULL,
    "conversation_id" bigint NULL,
    "content" text NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_messages_sender_id" FOREIGN KEY ("sender_id") REFERENCES "users"("id"),
    CONSTRAINT "fk_messages_conversation_id" FOREIGN KEY ("conversation_id") REFERENCES "conversations"("id")
);

-- Create index "idx_messages_deleted_at" to table: "messages"
CREATE INDEX "idx_messages_deleted_at" ON "public"."messages" ("deleted_at");

---- create above / drop below ----

DROP TABLE users CASCADE;

DROP TABLE conversations CASCADE;

DROP TABLE participants CASCADE;

DROP TABLE messages CASCADE;

