-- name: InsertTag :one

INSERT INTO tags ("name", "color_r", "color_g", "color_b", "color_a")
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateTag :exec
UPDATE tags
SET
    "name" = ?,
    "color_r" = ?,
    "color_g" = ?,
    "color_b" = ?,
    "color_a" = ?
WHERE "id" = ?;

-- name: DeleteTag :exec
DELETE FROM tags
WHERE "id" = ?;

-- name: GetAllTags :many
SELECT 
    *
FROM tags
ORDER BY "name";

-- name: UpsertTagSkill :one
INSERT INTO tag_skills ("tag_id", "skill_id", "skill_level")
VALUES (?, ?, ?)
ON CONFLICT ("tag_id", "skill_id") DO UPDATE
SET
    "skill_level" = excluded.skill_level
RETURNING *;

-- name: DeleteTagSkills :exec
DELETE FROM tag_skills
WHERE 
    "tag_id" = ?
    AND "skill_id" IN (sqlc.slice(skill_ids));

-- name: GetAllTagSkills :many
SELECT *
FROM tag_skills
WHERE "tag_id" = ?
ORDER BY "skill_id";