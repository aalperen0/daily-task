-- name: CreateTasks :one
INSERT INTO tasks(id, user_id, title, description, due_date, status, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetTaskOfUser :many
SELECT tasks.title, tasks.description, tasks.status, tasks.due_date
FROM tasks
JOIN users ON tasks.user_id = users.id
WHERE users.id = $1;


-- name: DeleteTaskOfUser :exec
DELETE FROM tasks
USING users
WHERE tasks.user_id = users.id AND users.id = $1;

-- name: DeleteTaskByTaskId :exec
DELETE FROM tasks
USING users
WHERE tasks.user_id = users.id AND tasks.id = $1;

-- name: UpdateTaskOfUser :exec
UPDATE tasks
SET description = $2,
    title = $3,
    due_date = $4,
    status = $5,
    updated_at = $6
WHERE tasks.id = $1;


-- name: GetTaskById :one
SELECT * FROM tasks
WHERE tasks.id = $1;