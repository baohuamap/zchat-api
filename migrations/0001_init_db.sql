
-- Create "users" table
CREATE TABLE "public"."users" (
    "id" bigserial NOT NULL,
    "created_at" timestamptz NULL,
    "updated_at" timestamptz NULL,
    "deleted_at" timestamptz NULL,
    "email" text NULL,
    "first" text NULL,
    "last" text NULL,
    "username" text NULL,
    "password" text NULL,
    PRIMARY KEY ("id")
);

-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");

-- Create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "public"."users" ("username");

-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");


-- Create "chats" table
CREATE TABLE "public"."chats" (
    "id" bigserial NOT NULL,
    "created_at" timestamptz NULL,
    "updated_at" timestamptz NULL,
    "deleted_at" timestamptz NULL,
    "sender_id" bigint NULL,
    "receiver_id" bigint NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_chats_sender_id" FOREIGN KEY ("sender_id") REFERENCES "users"("id"),
    CONSTRAINT "fk_chats_receiver_id" FOREIGN KEY ("receiver_id") REFERENCES "users"("id")
);

-- Create index "idx_files_deleted_at" to table: "files"
CREATE INDEX "idx_chats_deleted_at" ON "public"."chats" ("deleted_at");

-- Create "messages" table
CREATE TABLE "public"."messages" (
    "id" bigserial NOT NULL,
    "created_at" timestamptz NULL,
    "updated_at" timestamptz NULL,
    "deleted_at" timestamptz NULL,
    "chat_id" bigint NULL,
    "message" text NULL,
    "is_read" boolean NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_messages_chat_id" FOREIGN KEY ("chat_id") REFERENCES "chats"("id")
);

-- Create index "idx_messages_deleted_at" to table: "messages"
CREATE INDEX "idx_messages_deleted_at" ON "public"."messages" ("deleted_at");

---- create above / drop below ----

DROP TABLE users CASCADE;

