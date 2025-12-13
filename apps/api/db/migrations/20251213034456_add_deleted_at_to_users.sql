-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "deleted_at" timestamptz NULL;
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
