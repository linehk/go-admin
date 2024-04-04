-- name: GetUser :one
SELECT id, username, email, phone, remark, status, created, updated
FROM app_user
WHERE id = $1 LIMIT 1;

-- name: ListUser :many
SELECT id, username, email, phone, remark, status, created, updated
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
INSERT INTO app_user (username, password, email, phone, remark, status, created, updated)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateUser :exec
UPDATE app_user
SET username = $2, password = $3, email = $4, phone = $5, remark = $6, status = $7, created = $8, updated = $9
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM app_user
WHERE id = $1;