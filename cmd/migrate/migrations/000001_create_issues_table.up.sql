-- 000001_create_issues_table.up.sql
--
-- This is the "up" migration — it creates your database schema.
-- golang-migrate runs this file when you do `go run ./cmd/migrate -direction up`.
-- It tracks which migrations have been applied in a `schema_migrations` table
-- that it creates and manages automatically.

-- Create a custom PostgreSQL enum type.
-- PostgreSQL enums are stored efficiently (4 bytes) and enforce that only
-- these exact values can be inserted. This maps to your Go StatusType.
CREATE TYPE status_type AS ENUM ('Incomplete', 'In-Progress', 'Complete');

-- Create the issues table.
-- BIGSERIAL auto-generates sequential IDs (1, 2, 3, ...).
-- TIMESTAMPTZ stores timestamps with timezone info — always use this over
-- TIMESTAMP to avoid timezone bugs.
CREATE TABLE IF NOT EXISTS issues (
    id          BIGSERIAL PRIMARY KEY,
    author_id   UUID NOT NULL,
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status      status_type NOT NULL DEFAULT 'Incomplete',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
