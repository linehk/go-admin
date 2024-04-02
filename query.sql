-- name: GetUser :one
SELECT * FROM app_user
WHERE id = $1 LIMIT 1;

-- name: ListUser :many
SELECT * FROM app_user
ORDER BY username;

-- name: CreateUser :one
INSERT INTO app_user (
  username, password
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE app_user
  set username = $2,
  password = $3
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM app_user
WHERE id = $1;