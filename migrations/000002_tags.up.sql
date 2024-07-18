CREATE TABLE IF NOT EXISTS tags (
    "id" INTEGER PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "color_r" INTEGER NOT NULL,
    "color_g" INTEGER NOT NULL,
    "color_b" INTEGER NOT NULL
    -- , "color_a" INTEGER NOT NULL -- not really at this level, but for sqlc
);

CREATE TABLE IF NOT EXISTS tag_skills (
    "tag_id" BIGINT NOT NULL REFERENCES "tags" ("id") ON DELETE CASCADE,
    "skill_id" BIGINT NOT NULL,
    "skill_level" INTEGER NOT NULL,
    PRIMARY KEY ("tag_id", "skill_id")
);