-- Create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL,
  "name" character varying NOT NULL,
  "email" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
