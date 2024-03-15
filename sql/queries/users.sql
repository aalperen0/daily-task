-- name: CreateUser :one

INSERT INTO users(id, name, email, password, created_at)
VALUES($1, $2, $3, $4, $5)
RETURNING *;


-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;


-- name: DeleteUserById :exec
DELETE FROM users
WHERE id = $1;


-- name: UpdateUserById :exec
UPDATE users
SET name = $2,
email = $3,
password = $4
WHERE id = $1;


-- name: GetIdOfUser :one
SELECT id FROM users WHERE email = $1;

-- name: AuthenticateUser :many
SELECT email, password FROM users
WHERE email = $1 AND password = $2;