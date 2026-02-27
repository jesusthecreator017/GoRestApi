-- 000004_replace_authors_with_users.up.sql
--
-- Replace authors with users for issue ownership.
-- Truncates existing issues since author_id values don't map to users.

TRUNCATE TABLE issues;

ALTER TABLE issues
    DROP CONSTRAINT IF EXISTS fk_issues_author;

ALTER TABLE issues
    RENAME COLUMN author_id TO user_id;

ALTER TABLE issues
    ADD CONSTRAINT fk_issues_user
    FOREIGN KEY (user_id) REFERENCES users(id);

DROP TABLE IF EXISTS authors;

-- Insert dummy issues for development.
-- These reference the first user found in the users table.
INSERT INTO issues (title, user_id, description, status)
SELECT title, user_id, description, status::status_type
FROM (
    SELECT
        (SELECT id FROM users LIMIT 1) AS user_id,
        unnest(ARRAY[
            'Set up CI/CD pipeline',
            'Add unit tests for auth',
            'Fix login redirect bug',
            'Design dashboard layout',
            'Implement dark mode',
            'Write API documentation'
        ]) AS title,
        unnest(ARRAY[
            'Configure GitHub Actions for automated builds and deploys',
            'Cover register, login, and token validation with tests',
            'Users are not redirected back to the page they came from after login',
            'Create wireframes and mockups for the main dashboard view',
            'Add a theme toggle and persist preference in localStorage',
            'Document all REST endpoints with request/response examples'
        ]) AS description,
        unnest(ARRAY[
            'Incomplete',
            'Incomplete',
            'In-Progress',
            'In-Progress',
            'Complete',
            'Incomplete'
        ]) AS status
) AS seed
WHERE user_id IS NOT NULL;
