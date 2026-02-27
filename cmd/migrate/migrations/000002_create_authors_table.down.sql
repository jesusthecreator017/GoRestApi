-- 000002_create_authors_table.down.sql
--
-- Undo: drop the FK first, then the table.

ALTER TABLE issues DROP CONSTRAINT IF EXISTS fk_issues_author;
DROP TABLE IF EXISTS authors;
