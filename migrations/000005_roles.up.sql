DROP INDEX IF EXISTS "idx_tag_name_unique";

CREATE TABLE IF NOT EXISTS roles (
    "id" INTEGER PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "label" VARCHAR NOT NULL,
    "operator" VARCHAR NOT NULL,
    "color_r" INTEGER NOT NULL,
    "color_g" INTEGER NOT NULL,
    "color_b" INTEGER NOT NULL,
    "color_a" INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS role_tags (
    "role_id" INTEGER NOT NULL REFERENCES "roles" ("id") ON DELETE CASCADE,
    "tag_id" INTEGER NOT NULL REFERENCES "tags" ("id") ON DELETE CASCADE,
    PRIMARY KEY ("role_id", "tag_id")
);