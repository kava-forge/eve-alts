CREATE TABLE sqlite_master (
  "type" text,
  "name" text not null,
  "tbl_name" text,
  "rootpage" integer,
  "sql" text
);

CREATE TABLE IF NOT EXISTS trnTranslations (
        "tcID" INTEGER NOT NULL,
        "keyID" INTEGER NOT NULL,
        "languageID" VARCHAR(50) NOT NULL,
        "text" TEXT NOT NULL,
        PRIMARY KEY ("tcID", "keyID", "languageID")
);

CREATE INDEX IF NOT EXISTS "idx_translations_by_name" 
ON trnTranslations ("tcID", "languageID", LOWER("text"))
WHERE "tcID" = 8;