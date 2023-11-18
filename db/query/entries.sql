-- name: CreateEntry :one
INSERT INTO entries (
    account_id,
    amount
) VALUES ($1, $2)
RETURNING *;

-- name: GetEntry :one
SELECT *
FROM entries
WHERE id = $1 LIMIT 1;


-- name: GetEntriesByAccount :many
SELECT *
FROM entries
WHERE account_id = $1;


-- name: GetAllEntries :many
SELECT *
FROM entries
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateEntry :one
UPDATE entries
SET amount = $1,
    account_id = $2
WHERE id = $3
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries WHERE id = $1;
