-- name: CreateIssue :one
INSERT INTO issues (title, user_id, description, status)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, title, description, status, created_at, updated_at;

-- name: GetIssueByID :one
SELECT i.id, i.user_id, i.title, i.description, i.status, i.created_at, i.updated_at,
       u.name AS user_name
FROM issues i
JOIN users u ON u.id = i.user_id
WHERE i.id = $1;

-- name: ListIssues :many
SELECT i.id, i.user_id, i.title, i.description, i.status, i.created_at, i.updated_at,
       u.name AS user_name
FROM issues i
JOIN users u ON u.id = i.user_id
ORDER BY i.created_at DESC;

-- name: ListIssuesByUserID :many
SELECT i.id, i.user_id, i.title, i.description, i.status, i.created_at, i.updated_at,
       u.name AS user_name
FROM issues i
JOIN users u ON u.id = i.user_id
WHERE i.user_id = $1
ORDER BY i.created_at DESC;

-- name: UpdateIssueStatus :one
UPDATE issues
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING id, user_id, title, description, status, created_at, updated_at;

-- name: DeleteIssue :exec
DELETE FROM issues
WHERE id = $1;
