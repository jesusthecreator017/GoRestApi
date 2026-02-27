-- name: CreateUser :one
INSERT INTO users (email, password_hash, name)
VALUES ($1, $2, $3)
RETURNING id, email, name, password_hash, permissions, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT u.id, u.email, u.password_hash, u.name, u.permissions, u.created_at, u.updated_at
FROM users u
WHERE u.email = $1;

-- name: GetUserByID :one
SELECT u.id, u.email, u.password_hash, u.name, u.permissions, u.created_at, u.updated_at
FROM users u
WHERE u.id = $1;
