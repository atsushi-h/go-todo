-- Create "todos" table
CREATE TABLE "public"."todos" (
  "id" bigserial NOT NULL,
  "title" character varying NOT NULL,
  "description" character varying NULL,
  "completed" boolean NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
