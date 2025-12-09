-- Drop index "todos_user_id_idx" from table: "todos"
DROP INDEX "public"."todos_user_id_idx";
-- Modify "todos" table
ALTER TABLE "public"."todos" ALTER COLUMN "title" TYPE text, ALTER COLUMN "description" TYPE text, ALTER COLUMN "completed" SET NOT NULL, ALTER COLUMN "created_at" SET DEFAULT now(), ALTER COLUMN "updated_at" SET DEFAULT now(), ADD COLUMN "deleted_at" timestamptz NULL;
-- Create index "idx_todos_deleted_at" to table: "todos"
CREATE INDEX "idx_todos_deleted_at" ON "public"."todos" ("deleted_at");
-- Create index "idx_todos_user_id" to table: "todos"
CREATE INDEX "idx_todos_user_id" ON "public"."todos" ("user_id");
-- Modify "users" table
ALTER TABLE "public"."users" DROP CONSTRAINT "users_email_key", ALTER COLUMN "email" TYPE text, ALTER COLUMN "name" TYPE text, ALTER COLUMN "name" SET NOT NULL, ALTER COLUMN "avatar_url" TYPE text, ALTER COLUMN "provider" TYPE text, ALTER COLUMN "provider_id" TYPE text, ALTER COLUMN "created_at" SET DEFAULT now(), ALTER COLUMN "updated_at" SET DEFAULT now(), ADD CONSTRAINT "users_provider_provider_id_key" UNIQUE ("provider", "provider_id");
