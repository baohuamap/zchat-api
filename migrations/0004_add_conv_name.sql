ALTER TABLE "public"."conversations"
ADD COLUMN "name" text DEFAULT NULL;

---- create above / drop below ----

ALTER TABLE "public"."conversations"
DROP COLUMN "name";
