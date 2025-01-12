CREATE TYPE "users_role_enum" AS ENUM ('admin', 'cashier');

CREATE TABLE "users" (
    "id" BIGSERIAL PRIMARY KEY,
    "name" varchar NOT NULL,
    "email" varchar NUL NULL,
    "password" varchar NOT NULL,
    "role" users_role_enum DEFAULT 'cashier',
    "created_at" timestampz NOT NULL DEFAULT (now()),
    "updated_at" timestampz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX "email" ON "users" ("email");
