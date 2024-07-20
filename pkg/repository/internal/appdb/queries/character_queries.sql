-- name: UpsertCorporation :one
INSERT INTO corporations ("id", "name", "ticker", "alliance_id", "picture")
VALUES (?, ?, ?, ?, ?)
ON CONFLICT ("id") DO UPDATE
SET
    "name" = excluded.name,
    "picture" = excluded.picture,
    "alliance_id" = excluded.alliance_id,
    "ticker" = excluded.ticker
RETURNING *;

-- name: UpsertAlliance :one
INSERT INTO alliances ("id", "name", "ticker", "picture")
VALUES (?, ?, ?, ?)
ON CONFLICT ("id") DO UPDATE
SET
    "name" = excluded.name,
    "picture" = excluded.picture,
    "ticker" = excluded.ticker
RETURNING *;

-- name: UpsertCharacter :one
INSERT INTO characters ("id", "name", "picture", "corporation_id")
VALUES (?, ?, ?, ?)
ON CONFLICT ("id") DO UPDATE
SET
    "name" = excluded.name,
    "picture" = excluded.picture,
    "corporation_id" = excluded.corporation_id
RETURNING *;

-- name: DeleteCharacter :exec
DELETE FROM characters
WHERE "id" = ?;

-- name: UpsertToken :one

INSERT INTO tokens ("character_id", "access_token", "refresh_token", "token_type", "expiration")
VALUES (?, ?, ?, ?, ?)
ON CONFLICT ("character_id") DO UPDATE
SET
    "access_token" = excluded.access_token,
    "refresh_token" = excluded.refresh_token,
    "token_type" = excluded.token_type,
    "expiration" = excluded.expiration
RETURNING *;

-- name: GetAllCharacters :many
SELECT 
    sqlc.embed(characters),
    sqlc.embed(corporations),
    sqlc.embed(alliances)
FROM characters
INNER JOIN corporations ON characters."corporation_id" = corporations."id"
LEFT JOIN alliances ON corporations."alliance_id" = alliances."id";

-- name: GetTokenForCharacter :one
SELECT * 
FROM tokens
WHERE "character_id" = ?
LIMIT 1;

-- name: GetAllCharacterSkills :many
SELECT *
FROM character_skills
WHERE "character_id" = ?
ORDER BY "skill_id";

-- name: UpsertCharacterSkill :one
INSERT INTO character_skills ("character_id", "skill_id", "skill_level")
VALUES (?, ?, ?)
ON CONFLICT ("character_id", "skill_id") DO UPDATE
SET
    "skill_level" = excluded.skill_level
RETURNING *;

-- name: DeleteCharacterSkills :exec
DELETE FROM character_skills
WHERE 
    "character_id" = ?
    AND "skill_id" IN (sqlc.slice(skill_ids));