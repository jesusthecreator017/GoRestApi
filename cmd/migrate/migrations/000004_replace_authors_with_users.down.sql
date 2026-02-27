-- 000004_replace_authors_with_users.down.sql
--
-- Reverse: recreate authors table, rename user_id back to author_id.

CREATE TABLE IF NOT EXISTS authors (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

TRUNCATE TABLE issues;

ALTER TABLE issues
    DROP CONSTRAINT IF EXISTS fk_issues_user;

ALTER TABLE issues
    RENAME COLUMN user_id TO author_id;

ALTER TABLE issues
    ADD CONSTRAINT fk_issues_author
    FOREIGN KEY (author_id) REFERENCES authors(id);
