-- name: CountUsers :one
SELECT COUNT(*) AS total_users FROM users;

-- name: CountIssuesByStatus :many
SELECT status, COUNT(*) AS count FROM issues GROUP BY status;
