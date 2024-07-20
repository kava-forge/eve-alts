// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package appdb

import (
	"database/sql"
	"time"
)

type Alliance struct {
	ID      sql.NullInt64
	Name    sql.NullString
	Ticker  sql.NullString
	Picture sql.NullString
}

type Character struct {
	ID            int64
	Name          string
	Picture       string
	CorporationID int64
}

type CharacterSkill struct {
	CharacterID int64
	SkillID     int64
	SkillLevel  int64
}

type Corporation struct {
	ID         int64
	AllianceID sql.NullInt64
	Name       string
	Ticker     string
	Picture    string
}

type Role struct {
	ID       int64
	Name     string
	Label    string
	Operator string
	ColorR   int64
	ColorG   int64
	ColorB   int64
	ColorA   int64
}

type RoleTag struct {
	RoleID int64
	TagID  int64
}

type Tag struct {
	ID     int64
	Name   string
	ColorR int64
	ColorG int64
	ColorB int64
	ColorA int64
}

type TagSkill struct {
	TagID      int64
	SkillID    int64
	SkillLevel int64
}

type Token struct {
	ID           int64
	CharacterID  int64
	AccessToken  string
	RefreshToken string
	TokenType    string
	Expiration   time.Time
}
