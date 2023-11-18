-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    balance,
    currency
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccount :one
SELECT *
FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT *
FROM accounts
WHERE id = $1
LIMIT 1
FOR NO KEY UPDATE;

-- name: GetAccounts :many
SELECT *
FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET
    balance = CASE WHEN sqlc.arg(new_balance)::bigint != 0 THEN sqlc.arg(new_balance) ELSE balance END,
    owner = CASE WHEN sqlc.arg(new_owner)::varchar != '' THEN sqlc.arg(new_owner) ELSE owner END,
    currency = CASE WHEN sqlc.arg(new_currency)::varchar != '' THEN sqlc.arg(new_currency) ELSE currency END
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;
