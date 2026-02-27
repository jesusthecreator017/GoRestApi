-- 000001_create_issues_table.down.sql
--
-- This is the "down" migration — it undoes what the "up" migration did.
-- golang-migrate runs this when you do `go run ./cmd/migrate -direction down`.
-- Always drop in reverse order of creation: table first, then types it depends on.

DROP TABLE IF EXISTS issues;
DROP TYPE IF EXISTS status_type;
