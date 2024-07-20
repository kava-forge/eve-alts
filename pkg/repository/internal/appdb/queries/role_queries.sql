-- name: InsertRole :one

INSERT INTO roles ("name", "label", "operator", "color_r", "color_g", "color_b", "color_a")
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateRole :exec
UPDATE roles
SET
    "name" = ?,
    "label" = ?,
    "operator" = ?,
    "color_r" = ?,
    "color_g" = ?,
    "color_b" = ?,
    "color_a" = ?
WHERE "id" = ?;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE "id" = ?;

-- name: GetAllRoles :many
SELECT 
    *
FROM roles
ORDER BY "name";

-- name: UpsertRoleTag :one
INSERT INTO role_tags ("role_id", "tag_id")
VALUES (?, ?)
ON CONFLICT ("role_id", "tag_id") DO NOTHING
RETURNING *;

-- name: DeleteRoleTags :exec
DELETE FROM role_tags
WHERE 
    "role_id" = ?
    AND "tag_id" IN (sqlc.slice(tag_ids));

-- name: GetAllRoleTags :many
SELECT tags.*
FROM tags
JOIN role_tags ON tags."id" = role_tags."tag_id"
WHERE role_tags."role_id" = ?
ORDER BY tags."name";