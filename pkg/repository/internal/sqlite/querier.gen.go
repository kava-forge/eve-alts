// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package sqlite

import (
	"context"
)

type Querier interface {
	DeleteCharacter(ctx context.Context, db DBTX, id int64) error
	DeleteCharacterSkills(ctx context.Context, db DBTX, arg DeleteCharacterSkillsParams) error
	DeleteTag(ctx context.Context, db DBTX, id int64) error
	DeleteTagSkills(ctx context.Context, db DBTX, arg DeleteTagSkillsParams) error
	GetAllCharacterSkills(ctx context.Context, db DBTX, characterID int64) ([]CharacterSkill, error)
	GetAllCharacters(ctx context.Context, db DBTX) ([]GetAllCharactersRow, error)
	GetAllTagSkills(ctx context.Context, db DBTX, tagID int64) ([]TagSkill, error)
	GetAllTags(ctx context.Context, db DBTX) ([]Tag, error)
	GetTokenForCharacter(ctx context.Context, db DBTX, characterID int64) (Token, error)
	InsertTag(ctx context.Context, db DBTX, arg InsertTagParams) (Tag, error)
	UpdateTag(ctx context.Context, db DBTX, arg UpdateTagParams) error
	UpsertAlliance(ctx context.Context, db DBTX, arg UpsertAllianceParams) (Alliance, error)
	UpsertCharacter(ctx context.Context, db DBTX, arg UpsertCharacterParams) (Character, error)
	UpsertCharacterSkill(ctx context.Context, db DBTX, arg UpsertCharacterSkillParams) (CharacterSkill, error)
	UpsertCorporation(ctx context.Context, db DBTX, arg UpsertCorporationParams) (Corporation, error)
	UpsertTagSkill(ctx context.Context, db DBTX, arg UpsertTagSkillParams) (TagSkill, error)
	UpsertToken(ctx context.Context, db DBTX, arg UpsertTokenParams) (Token, error)
}

var _ Querier = (*Queries)(nil)
