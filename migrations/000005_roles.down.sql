DROP TABLE IF EXISTS role_tags;

DROP TABLE IF EXISTS roles;

CREATE UNIQUE INDEX IF NOT EXISTS "idx_tag_name_unique" ON tags ("name");