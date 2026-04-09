-- name: CreateRecipe :one
INSERT INTO recipe (
    name, ingredients, cost
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: UpdateRecipe :one
UPDATE recipe
SET name = $2,
    ingredients = $3,
    cost = $4
WHERE id = $1
RETURNING *;

-- name: GetRecipe :one
SELECT * FROM recipe
WHERE id = $1;

-- name: ListRecipe :many
SELECT id, name, ingredients, cost
FROM recipe
LIMIT $1 OFFSET $2;

-- name: DeleteRecipe :exec
DELETE FROM recipe
WHERE id = $1;
