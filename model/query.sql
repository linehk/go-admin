--------------------------------- User --------------------------------
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
SELECT EXISTS (SELECT 1 FROM app_user WHERE id = $1);

-- name: CheckUserByUsername :one
SELECT EXISTS (SELECT 1 FROM app_user WHERE username = $1);

-- name: CreateUser :one
INSERT INTO app_user (username, password, name, email, phone, remark, status,
created, updated)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateUser :one
UPDATE app_user
SET username = $2, password = $3, name = $4, email = $5, phone = $6,
remark = $7, status = $8, created = $9, updated = $10
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM app_user
WHERE id = $1;


--------------------------------- UserRole --------------------------------
-- name: GetUserRole :one
SELECT *
FROM user_role
WHERE id = $1 LIMIT 1;

-- name: ListUserRoleByUserIDList :many
SELECT *
FROM user_role
WHERE user_id = ANY($1::int[]);

-- name: CheckUserRoleByID :one
SELECT EXISTS (SELECT 1 FROM user_role WHERE id = $1);

-- name: CreateUserRole :one
INSERT INTO user_role (user_id, role_id, created, updated)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUserRole :one
UPDATE user_role
SET user_id = $2, role_id = $3, created = $4, updated = $5
WHERE id = $1
RETURNING *;

-- name: DeleteUserRole :exec
DELETE FROM user_role
WHERE id = $1;

-- name: DeleteUserRoleByUserID :exec
DELETE FROM user_role
WHERE user_id = $1;


--------------------------------- Role --------------------------------
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
SELECT EXISTS (SELECT 1 FROM role WHERE id = $1);

-- name: CheckRoleByCode :one
SELECT EXISTS (SELECT 1 FROM role WHERE code = $1);

-- name: CreateRole :one
INSERT INTO role (code, name, description, sequence, status, created, updated)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateRole :one
UPDATE role
SET code = $2, name = $3, description = $4, sequence = $5, status = $6,
created = $7, updated = $8
WHERE id = $1
RETURNING *;

-- name: DeleteRole :exec
DELETE FROM role
WHERE id = $1;


--------------------------------- RoleMenu --------------------------------
-- name: GetRoleMenu :one
SELECT *
FROM role_menu
WHERE id = $1 LIMIT 1;

-- name: ListRoleMenuByRoleIDList :many
SELECT *
FROM role_menu
WHERE role_id = ANY($1::int[]);

-- name: CheckRoleMenuByID :one
SELECT EXISTS (SELECT 1 FROM role_menu WHERE id = $1);

-- name: CreateRoleMenu :one
INSERT INTO role_menu (role_id, menu_id, created, updated)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateRoleMenu :one
UPDATE role_menu
SET role_id = $2, menu_id = $3, created = $4, updated = $5
WHERE id = $1
RETURNING *;

-- name: DeleteRoleMenu :exec
DELETE FROM role_menu
WHERE id = $1;


--------------------------------- Menu --------------------------------
-- name: GetMenu :one
SELECT *
FROM menu
WHERE id = $1 LIMIT 1;

-- name: CheckMenuByID :one
SELECT EXISTS (SELECT 1 FROM menu WHERE id = $1);

-- name: CreateMenu :one
INSERT INTO menu (code, name, description, sequence, type, path, property,
parent_id, parent_path, status, created, updated)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: UpdateMenu :one
UPDATE menu
SET code = $2, name = $3, description = $4, sequence = $5, type = $6,
path = $7, property = $8, parent_id = $9, parent_path = $10, status = $11,
created = $12, updated = $13
WHERE id = $1
RETURNING *;

-- name: DeleteMenu :exec
DELETE FROM menu
WHERE id = $1;


--------------------------------- Resource --------------------------------
-- name: GetResource :one
SELECT *
FROM resource
WHERE id = $1 LIMIT 1;

-- name: ListResourceByMenuIDList :many
SELECT *
FROM resource
WHERE menu_id = ANY($1::int[]);

-- name: CheckResourceByID :one
SELECT EXISTS (SELECT 1 FROM resource WHERE id = $1);

-- name: CreateResource :one
INSERT INTO resource (menu_id, method, path, created, updated)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateResource :one
UPDATE resource
SET menu_id = $2, method = $3, path = $4, created = $5, updated = $6
WHERE id = $1
RETURNING *;

-- name: DeleteResource :exec
DELETE FROM resource
WHERE id = $1;