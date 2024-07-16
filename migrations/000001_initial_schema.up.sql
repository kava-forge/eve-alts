CREATE TABLE IF NOT EXISTS alliances (
    "id" BIGINT UNIQUE,
    "name" VARCHAR,
    "ticker" VARCHAR,
    "picture" VARCHAR
);

CREATE TABLE IF NOT EXISTS corporations (
    "id" BIGINT NOT NULL PRIMARY KEY,
    "alliance_id" BIGINT REFERENCES "alliances" ("id")  ON DELETE CASCADE,
    "name" VARCHAR NOT NULL,
    "ticker" VARCHAR NOT NULL,
    "picture" VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS characters (
    "id" BIGINT NOT NULL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "picture" VARCHAR NOT NULL,
    "corporation_id" BIGINT NOT NULL REFERENCES "corporations" ("id")  ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tokens (
    "id" INTEGER PRIMARY KEY,
    "character_id" BIGINT NOT NULL REFERENCES "characters" ("id") ON DELETE CASCADE UNIQUE,
    "access_token" TEXT NOT NULL,
    "refresh_token" TEXT NOT NULL,
    "token_type" VARCHAR NOT NULL,
    "expiration" TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS character_skills (
    "character_id" BIGINT NOT NULL REFERENCES "characters" ("id") ON DELETE CASCADE,
    "skill_id" BIGINT NOT NULL,
    "skill_level" INTEGER NOT NULL,
    PRIMARY KEY ("character_id", "skill_id")
);