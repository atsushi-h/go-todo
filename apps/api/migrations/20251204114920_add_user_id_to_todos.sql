-- 既存データを削除
TRUNCATE TABLE "public"."todos";

-- Modify "todos" table
ALTER TABLE "public"."todos" ADD COLUMN "user_id" bigint NOT NULL;

-- 外部キー制約を追加
ALTER TABLE "public"."todos"
  ADD CONSTRAINT "todos_user_id_fkey"
  FOREIGN KEY ("user_id")
  REFERENCES "public"."users" ("id")
  ON DELETE CASCADE;

-- インデックスを追加
CREATE INDEX "todos_user_id_idx" ON "public"."todos" ("user_id");
