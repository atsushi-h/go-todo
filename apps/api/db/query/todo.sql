-- name: GetTodoByID :one
SELECT * FROM todos
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: ListTodosByUser :many
SELECT * FROM todos
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CreateTodo :one
INSERT INTO todos (user_id, title, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateTodo :one
UPDATE todos
SET
    title = COALESCE(sqlc.narg(title), title),
    description = COALESCE(sqlc.narg(description), description),
    completed = COALESCE(sqlc.narg(completed), completed),
    updated_at = NOW()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteTodo :exec
UPDATE todos
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;
