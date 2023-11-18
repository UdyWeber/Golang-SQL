-- name: CreatePost :one
INSERT INTO posts (from_account_id, to_account_id, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPost :one
SELECT *
FROM posts
WHERE id = $1
LIMIT 1;

-- name: GetPosts :many
SELECT *
FROM posts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdatePost :one
UPDATE posts
SET
    amount = $1,
    to_account_id = $2,
    from_account_id = $3
WHERE
    id = $4
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;