// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: tag_queries.sql

package appdb

import (
	"context"
	"strings"
)

const deleteTag = `-- name: DeleteTag :exec
DELETE FROM tags
WHERE "id" = ?
`

func (q *Queries) DeleteTag(ctx context.Context, db DBTX, id int64) error {
	_, err := db.ExecContext(ctx, deleteTag, id)
	return err
}

const deleteTagSkills = `-- name: DeleteTagSkills :exec
DELETE FROM tag_skills
WHERE 
    "tag_id" = ?
    AND "skill_id" IN (/*SLICE:skill_ids*/?)
`

type DeleteTagSkillsParams struct {
	TagID    int64
	SkillIds []int64
}

func (q *Queries) DeleteTagSkills(ctx context.Context, db DBTX, arg DeleteTagSkillsParams) error {
	query := deleteTagSkills
	var queryParams []interface{}
	queryParams = append(queryParams, arg.TagID)
	if len(arg.SkillIds) > 0 {
		for _, v := range arg.SkillIds {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:skill_ids*/?", strings.Repeat(",?", len(arg.SkillIds))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:skill_ids*/?", "NULL", 1)
	}
	_, err := db.ExecContext(ctx, query, queryParams...)
	return err
}

const getAllTagSkills = `-- name: GetAllTagSkills :many
SELECT tag_id, skill_id, skill_level
FROM tag_skills
WHERE "tag_id" = ?
ORDER BY "skill_id"
`

func (q *Queries) GetAllTagSkills(ctx context.Context, db DBTX, tagID int64) ([]TagSkill, error) {
	rows, err := db.QueryContext(ctx, getAllTagSkills, tagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TagSkill
	for rows.Next() {
		var i TagSkill
		if err := rows.Scan(&i.TagID, &i.SkillID, &i.SkillLevel); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllTags = `-- name: GetAllTags :many
SELECT 
    id, name, color_r, color_g, color_b, color_a
FROM tags
ORDER BY "name"
`

func (q *Queries) GetAllTags(ctx context.Context, db DBTX) ([]Tag, error) {
	rows, err := db.QueryContext(ctx, getAllTags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Tag
	for rows.Next() {
		var i Tag
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ColorR,
			&i.ColorG,
			&i.ColorB,
			&i.ColorA,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertTag = `-- name: InsertTag :one

INSERT INTO tags ("name", "color_r", "color_g", "color_b", "color_a")
VALUES (?, ?, ?, ?, ?)
RETURNING id, name, color_r, color_g, color_b, color_a
`

type InsertTagParams struct {
	Name   string
	ColorR int64
	ColorG int64
	ColorB int64
	ColorA int64
}

func (q *Queries) InsertTag(ctx context.Context, db DBTX, arg InsertTagParams) (Tag, error) {
	row := db.QueryRowContext(ctx, insertTag,
		arg.Name,
		arg.ColorR,
		arg.ColorG,
		arg.ColorB,
		arg.ColorA,
	)
	var i Tag
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ColorR,
		&i.ColorG,
		&i.ColorB,
		&i.ColorA,
	)
	return i, err
}

const updateTag = `-- name: UpdateTag :exec
UPDATE tags
SET
    "name" = ?,
    "color_r" = ?,
    "color_g" = ?,
    "color_b" = ?,
    "color_a" = ?
WHERE "id" = ?
`

type UpdateTagParams struct {
	Name   string
	ColorR int64
	ColorG int64
	ColorB int64
	ColorA int64
	ID     int64
}

func (q *Queries) UpdateTag(ctx context.Context, db DBTX, arg UpdateTagParams) error {
	_, err := db.ExecContext(ctx, updateTag,
		arg.Name,
		arg.ColorR,
		arg.ColorG,
		arg.ColorB,
		arg.ColorA,
		arg.ID,
	)
	return err
}

const upsertTagSkill = `-- name: UpsertTagSkill :one
INSERT INTO tag_skills ("tag_id", "skill_id", "skill_level")
VALUES (?, ?, ?)
ON CONFLICT ("tag_id", "skill_id") DO UPDATE
SET
    "skill_level" = excluded.skill_level
RETURNING tag_id, skill_id, skill_level
`

type UpsertTagSkillParams struct {
	TagID      int64
	SkillID    int64
	SkillLevel int64
}

func (q *Queries) UpsertTagSkill(ctx context.Context, db DBTX, arg UpsertTagSkillParams) (TagSkill, error) {
	row := db.QueryRowContext(ctx, upsertTagSkill, arg.TagID, arg.SkillID, arg.SkillLevel)
	var i TagSkill
	err := row.Scan(&i.TagID, &i.SkillID, &i.SkillLevel)
	return i, err
}
