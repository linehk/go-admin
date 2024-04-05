-- name: GetUser :one
SELECT username, name, email, phone, remark, status, created, updated
FROM app_user
WHERE id = $1 LIMIT 1;

-- name: ListUser :many
SELECT id, username, name, email, phone, remark, status, created, updated
FROM app_user
WHERE username LIKE $1 AND name LIKE $2 AND status = $3
ORDER BY created DESC;

-- name: CheckUserByID :one
SELECT 1
FROM app_user
WHERE id = $1 LIMIT 1;

-- name: CheckUserByUsername :one
SELECT 1
FROM app_user
WHERE username = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO app_user (username, password, name, email, phone, remark, status, created, updated)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING username, name, email, phone, remark, status, created, updated;

-- name: UpdateUser :exec
UPDATE app_user
SET username = $2, password = $3, name = $4, email = $5, phone = $6, remark = $7, status = $8, created = $9, updated = $10
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM app_user
WHERE id = $1;