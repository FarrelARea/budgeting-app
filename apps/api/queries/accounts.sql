-- name: ListAccounts :many
SELECT id, name, currency, balance, created_at, updated_at
FROM accounts
ORDER BY created_at DESC;

-- name: GetAccount :one
SELECT id, name, currency, balance, created_at, updated_at
FROM accounts
WHERE id = $1;

-- name: CreateAccount :one
INSERT INTO accounts (name, currency, balance)
VALUES ($1, $2, 0)
RETURNING id, name, currency, balance, created_at, updated_at;
