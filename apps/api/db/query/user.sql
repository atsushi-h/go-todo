-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByProviderID :one
SELECT * FROM users
WHERE provider = $1 AND provider_id = $2 AND deleted_at IS NULL;

-- name: CreateUser :one
INSERT INTO users (email, name, avatar_url, provider, provider_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    name = $2,
    avatar_url = $3,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteTodosByUserID :exec
UPDATE todos
SET deleted_at = NOW(), updated_at = NOW()
WHERE user_id = $1 AND deleted_at IS NULL;
