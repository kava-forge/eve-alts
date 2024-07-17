-- name: GetTableNames :many
SELECT 
    "name"
FROM sqlite_master
WHERE
    "type" = 'table'
ORDER BY "name";

-- name: CleanTranslations :exec
DELETE FROM trnTranslations
WHERE "tcID" != 8;