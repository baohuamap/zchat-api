ALTER TABLE "public"."conversations"
ADD COLUMN seen BOOLEAN DEFAULT FALSE;

---- create above / drop below ----

ALTER TABLE "public"."conversations"
DROP COLUMN seen;
