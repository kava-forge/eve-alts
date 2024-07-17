-- name: GetSkillName :one
SELECT
    "text" as skill_name
FROM
    trnTranslations
WHERE
    "tcID" = 8
    AND "keyID" = sqlc.arg(skill_id)
    AND "languageID" = sqlc.arg(language)
;

-- name: BatchGetSkillNames :many
SELECT
    "keyID" as skill_id,
    "text" as skill_name
FROM
    trnTranslations
WHERE
    "tcID" = 8
    AND "languageID" = sqlc.arg(language)
    AND "keyID" IN (sqlc.slice(skill_ids))
;

-- name: GetSkillIDFromName :one
SELECT
    "keyID" as skill_id
FROM
    trnTranslations
WHERE
    "tcID" = 8
    AND "languageID" = sqlc.arg(language)
    AND LOWER("text") = sqlc.arg(skill_name_lower)
LIMIT 1
;