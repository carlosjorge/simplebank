-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account,
    to_account,
    amount
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateTransfer :one
UPDATE transfers
SET from_account = $2,
    to_account = $3,
    amount = $4
WHERE id = $1
RETURNING *;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;
