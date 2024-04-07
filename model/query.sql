-- User
-- name: GetUser :one
SELECT *
FROM app_user
WHERE id = $1 LIMIT 1;

-- name: ListUser :many
SELECT *
FROM app_user
WHERE ($1::VARCHAR = '' OR $1::VARCHAR ILIKE '%' || $1 || '%')
AND ($2::VARCHAR = '' OR $2::VARCHAR ILIKE '%' || $2 || '%')
AND ($3::VARCHAR = '' OR $3::VARCHAR = $3)
AND id > $4
ORDER BY created DESC
LIMIT $5;

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
RETURNING *;

-- name: UpdateUser :one
UPDATE app_user
SET username = $2, password = $3, name = $4, email = $5, phone = $6, remark = $7, status = $8, created = $9, updated = $10
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM app_user
WHERE id = $1;

-- Role
-- name: GetRole :one
SELECT *
FROM role
WHERE id = $1 LIMIT 1;

-- name: ListRole :many
SELECT *
FROM role
WHERE ($1::VARCHAR = '' OR $1::VARCHAR ILIKE '%' || $1 || '%')
AND ($2::VARCHAR = '' OR $2::VARCHAR = $2)
AND id > $3
ORDER BY sequence, created DESC
LIMIT $4;

-- name: CheckRoleByID :one
SELECT 1
FROM role
WHERE id = $1 LIMIT 1;

-- name: CheckRoleByCode :one
SELECT 1
FROM role
WHERE code = $1 LIMIT 1;

-- name: CreateRole :one
INSERT INTO role (code, name, description, sequence, status, created, updated)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateRole :one
UPDATE role
SET code = $2, name = $3, description = $4, sequence = $5, status = $6, created = $7, updated = $8
WHERE id = $1
RETURNING *;

-- name: DeleteRole :exec
DELETE FROM role
WHERE id = $1;