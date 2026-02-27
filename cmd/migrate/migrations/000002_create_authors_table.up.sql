-- 000002_create_authors_table.up.sql
--
-- Creates the authors table and adds a foreign key from issues.author_id.

CREATE TABLE IF NOT EXISTS authors (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE issues
    ADD CONSTRAINT fk_issues_author
    FOREIGN KEY (author_id) REFERENCES authors(id);
