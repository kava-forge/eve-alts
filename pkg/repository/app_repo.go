package repository

import (
	"context"
	"database/sql"
	"image/color"
	"strconv"
	"time"

	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"

	"github.com/kava-forge/eve-alts/pkg/database"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/repository/internal/appdb"
	"github.com/kava-forge/eve-alts/pkg/telemetry"
)

type (
	Character      = appdb.Character
	Corporation    = appdb.Corporation
	Alliance       = appdb.Alliance
	Token          = appdb.Token
	CharacterSkill = appdb.CharacterSkill
	Tag            = appdb.Tag
	TagSkill       = appdb.TagSkill
)

type CharacterDBData struct {
	Character   Character
	Corporation Corporation
	Alliance    Alliance
	Skills      []CharacterSkill
}

type TagDBData struct {
	Tag    Tag
	Skills []TagSkill
}

func (t TagDBData) Color() color.Color {
	return color.RGBA{
		R: uint8(t.Tag.ColorR),
		G: uint8(t.Tag.ColorG),
		B: uint8(t.Tag.ColorB),
		A: uint8(t.Tag.ColorA),
	}
}

func (t TagDBData) StrID() string {
	return strconv.FormatInt(t.Tag.ID, 10)
}

//counterfeiter:generate . AppData
type AppData interface {
	UpsertCharacter(ctx context.Context, charID int64, name, picture string, corporationID int64, tx database.Tx) (Character, error)
	UpsertCorporation(ctx context.Context, corpID int64, name, ticker, picture string, allianceID int64, tx database.Tx) (Corporation, error)
	UpsertAlliance(ctx context.Context, allianceID int64, name, ticker, picture string, tx database.Tx) (Alliance, error)
	UpsertToken(ctx context.Context, charID int64, accessToken, refreshToken, tokenType string, expiration time.Time, tx database.Tx) (Token, error)
	GetAllCharacters(ctx context.Context, tx database.Tx) ([]*CharacterDBData, error)
	GetTokenForCharacter(ctx context.Context, charID int64, tx database.Tx) (Token, error)
	GetAllCharacterSkills(ctx context.Context, charID int64, tx database.Tx) ([]CharacterSkill, error)
	UpsertCharacterSkill(ctx context.Context, charID, skillID, trainedLevel int64, tx database.Tx) (CharacterSkill, error)
	DeleteCharacterSkills(ctx context.Context, charID int64, skillIDs []int64, tx database.Tx) error
	DeleteCharacter(ctx context.Context, charID int64, tx database.Tx) error
	InsertTag(ctx context.Context, name string, c color.Color, tx database.Tx) (Tag, error)
	UpdateTag(ctx context.Context, tagID int64, name string, c color.Color, tx database.Tx) error
	DeleteTag(ctx context.Context, tagID int64, tx database.Tx) error
	GetAllTags(ctx context.Context, tx database.Tx) ([]*TagDBData, error)
	GetAllTagSkills(ctx context.Context, tagID int64, tx database.Tx) ([]TagSkill, error)
	UpsertTagSkill(ctx context.Context, tagID, skillID, skillLevel int64, tx database.Tx) (TagSkill, error)
	DeleteTagSkills(ctx context.Context, tagID int64, skillIDs []int64, tx database.Tx) error
}

type appDependencies interface {
	DB() database.Connection
	Logger() logging.Logger
	Telemetry() *telemetry.Telemeter
}

type AppSqliteRepository struct {
	deps    appDependencies
	queries *appdb.Queries
}

var _ AppData = (*AppSqliteRepository)(nil)

func NewAppData(deps appDependencies) *AppSqliteRepository {
	return &AppSqliteRepository{
		deps:    deps,
		queries: appdb.New(),
	}
}

func (r *AppSqliteRepository) db(tx database.Tx) appdb.DBTX {
	if tx != nil {
		return tx
	}
	return r.deps.DB()
}

func (r *AppSqliteRepository) UpsertCharacter(ctx context.Context, charID int64, name, picture string, corporationID int64, tx database.Tx) (char Character, err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "UpsertCharacter")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling UpsertCharacter", keys.CharacterID, charID, keys.CharacterName, name, keys.CharacterPicture, picture, keys.CorporationID, corporationID)

	inner := func(ctx context.Context, tx database.Tx) error {
		char, err = r.queries.UpsertCharacter(ctx, tx, appdb.UpsertCharacterParams{
			ID:            charID,
			Name:          name,
			Picture:       picture,
			CorporationID: corporationID,
		})
		return errors.Wrap(err, "could not UpsertCharacter")
	}

	if tx == nil {
		err = errors.Wrap(database.TransactWithRetries(ctx, r.deps.Telemetry(), r.deps.Logger(), r.deps.DB(), &sql.TxOptions{}, inner), "could not TransactWithRetries")
	} else {
		err = inner(ctx, tx)
	}
	return char, err
}

func (r *AppSqliteRepository) UpsertCorporation(ctx context.Context, corpID int64, name, ticker, picture string, allianceID int64, tx database.Tx) (corp Corporation, err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "UpsertCharacter")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling UpsertCorporation", keys.CorporationID, corpID, keys.CorportaionName, name, keys.CorporationPicture, picture, keys.CorporationTicker, ticker, keys.AllianceID, allianceID)

	var allyID sql.NullInt64
	if allianceID != 0 {
		allyID.Int64 = allianceID
		allyID.Valid = true
	}

	inner := func(ctx context.Context, tx database.Tx) error {
		corp, err = r.queries.UpsertCorporation(ctx, tx, appdb.UpsertCorporationParams{
			ID:         corpID,
			Name:       name,
			Ticker:     ticker,
			Picture:    picture,
			AllianceID: allyID,
		})
		return errors.Wrap(err, "could not UpsertCorporation")
	}

	if tx == nil {
		err = errors.Wrap(database.TransactWithRetries(ctx, r.deps.Telemetry(), r.deps.Logger(), r.deps.DB(), &sql.TxOptions{}, inner), "could not TransactWithRetries")
	} else {
		err = inner(ctx, tx)
	}
	return corp, err
}

func (r *AppSqliteRepository) UpsertAlliance(ctx context.Context, allyID int64, name, ticker, picture string, tx database.Tx) (ally Alliance, err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "UpsertAlliance")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling UpsertAlliance", keys.AllianceID, allyID, keys.AllianceName, name, keys.AlliancePicture, picture, keys.AllianceTicker, ticker)

	inner := func(ctx context.Context, tx database.Tx) error {
		ally, err = r.queries.UpsertAlliance(ctx, tx, appdb.UpsertAllianceParams{
			ID:      sql.NullInt64{Int64: allyID, Valid: true},
			Name:    sql.NullString{String: name, Valid: true},
			Ticker:  sql.NullString{String: ticker, Valid: true},
			Picture: sql.NullString{String: picture, Valid: true},
		})
		return errors.Wrap(err, "could not UpsertAlliance")
	}

	if tx == nil {
		err = errors.Wrap(database.TransactWithRetries(ctx, r.deps.Telemetry(), r.deps.Logger(), r.deps.DB(), &sql.TxOptions{}, inner), "could not TransactWithRetries")
	} else {
		err = inner(ctx, tx)
	}
	return ally, err
}

func (r *AppSqliteRepository) UpsertToken(ctx context.Context, charID int64, accessToken, refreshToken, tokenType string, expiration time.Time, tx database.Tx) (tok Token, err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "UpsertToken")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling UpsertToken", keys.CharacterID, charID)

	inner := func(ctx context.Context, tx database.Tx) error {
		tok, err = r.queries.UpsertToken(ctx, tx, appdb.UpsertTokenParams{
			CharacterID:  charID,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    tokenType,
			Expiration:   expiration,
		})
		return errors.Wrap(err, "could not UpsertCharacter")
	}

	if tx == nil {
		err = errors.Wrap(database.TransactWithRetries(ctx, r.deps.Telemetry(), r.deps.Logger(), r.deps.DB(), &sql.TxOptions{}, inner), "could not TransactWithRetries")
	} else {
		err = inner(ctx, tx)
	}
	return tok, err
}

func (r *AppSqliteRepository) GetAllCharacters(ctx context.Context, tx database.Tx) (_ []*CharacterDBData, err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "GetAllCharacters")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling GetAllCharacters")

	chars, err := r.queries.GetAllCharacters(ctx, r.db(tx))
	if err != nil {
		return nil, errors.Wrap(err, "could not GetAllCharacters")
	}

	charDBData := make([]*CharacterDBData, 0, len(chars))
	for _, c := range chars {
		skills, err := r.GetAllCharacterSkills(ctx, c.Character.ID, tx)
		if err != nil && !errors.Is(err, database.ErrNoRows) {
			return nil, errors.Wrap(err, "could not GetAllCharacterSkills")
		}

		charDBData = append(charDBData, &CharacterDBData{
			Character:   c.Character,
			Corporation: c.Corporation,
			Alliance:    c.Alliance,
			Skills:      skills,
		})
	}

	return charDBData, nil
}

func (r *AppSqliteRepository) GetTokenForCharacter(ctx context.Context, charID int64, tx database.Tx) (_ Token, err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "GetTokenForCharacter")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling GetTokenForCharacter")

	tok, err := r.queries.GetTokenForCharacter(ctx, r.db(tx), charID)
	if err != nil {
		return tok, err
	}

	return tok, nil
}

func (r *AppSqliteRepository) GetAllCharacterSkills(ctx context.Context, charID int64, tx database.Tx) (_ []CharacterSkill, err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "GetAllCharacterSkills")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling GetAllCharacterSkills")

	skills, err := r.queries.GetAllCharacterSkills(ctx, r.db(tx), charID)
	if err != nil {
		return skills, err
	}

	return skills, nil
}

func (r *AppSqliteRepository) UpsertCharacterSkill(ctx context.Context, charID, skillID, skillLevel int64, tx database.Tx) (skill CharacterSkill, err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "UpsertCharacterSkill")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling UpsertCharacterSkill", keys.CharacterID, charID)

	inner := func(ctx context.Context, tx database.Tx) error {
		skill, err = r.queries.UpsertCharacterSkill(ctx, tx, appdb.UpsertCharacterSkillParams{
			CharacterID: charID,
			SkillID:     skillID,
			SkillLevel:  skillLevel,
		})
		return errors.Wrap(err, "could not UpsertCharacterSkill")
	}

	if tx == nil {
		err = errors.Wrap(database.TransactWithRetries(ctx, r.deps.Telemetry(), r.deps.Logger(), r.deps.DB(), &sql.TxOptions{}, inner), "could not TransactWithRetries")
	} else {
		err = inner(ctx, tx)
	}
	return skill, err
}

func (r *AppSqliteRepository) DeleteCharacterSkills(ctx context.Context, charID int64, skillIDs []int64, tx database.Tx) (err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "DeleteCharacterSkills")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling DeleteCharacterSkills", keys.CharacterID, charID, keys.SkillID, skillIDs)

	inner := func(ctx context.Context, tx database.Tx) error {
		err = r.queries.DeleteCharacterSkills(ctx, tx, appdb.DeleteCharacterSkillsParams{
			CharacterID: charID,
			SkillIds:    skillIDs,
		})
		return errors.Wrap(err, "could not UpsertCharacterSkill")
	}

	if tx == nil {
		err = errors.Wrap(database.TransactWithRetries(ctx, r.deps.Telemetry(), r.deps.Logger(), r.deps.DB(), &sql.TxOptions{}, inner), "could not TransactWithRetries")
	} else {
		err = inner(ctx, tx)
	}
	return err
}

func (r *AppSqliteRepository) DeleteCharacter(ctx context.Context, charID int64, tx database.Tx) (err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "DeleteCharacter")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling DeleteCharacter", keys.CharacterID, charID)

	inner := func(ctx context.Context, tx database.Tx) error {
		err = r.queries.DeleteCharacter(ctx, tx, charID)
		return errors.Wrap(err, "could not DeleteCharacter")
	}

	if tx == nil {
		err = errors.Wrap(database.TransactWithRetries(ctx, r.deps.Telemetry(), r.deps.Logger(), r.deps.DB(), &sql.TxOptions{}, inner), "could not TransactWithRetries")
	} else {
		err = inner(ctx, tx)
	}
	return err
}

func (r *AppSqliteRepository) InsertTag(ctx context.Context, name string, c color.Color, tx database.Tx) (tag Tag, err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "InsertTag")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling InsertTag", keys.TagName, name)

	cr, cg, cb, ca := c.RGBA()

	inner := func(ctx context.Context, tx database.Tx) error {
		tag, err = r.queries.InsertTag(ctx, tx, appdb.InsertTagParams{
			Name:   name,
			ColorR: int64(cr),
			ColorG: int64(cg),
			ColorB: int64(cb),
			ColorA: int64(ca),
		})
		return errors.Wrap(err, "could not InsertTag")
	}

	if tx == nil {
		err = errors.Wrap(database.TransactWithRetries(ctx, r.deps.Telemetry(), r.deps.Logger(), r.deps.DB(), &sql.TxOptions{}, inner), "could not TransactWithRetries")
	} else {
		err = inner(ctx, tx)
	}
	return tag, err
}

func (r *AppSqliteRepository) UpdateTag(ctx context.Context, tagID int64, name string, c color.Color, tx database.Tx) (err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "UpdateTag")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling UpdateTag", keys.TagID, tagID, keys.TagName, name)

	cr, cg, cb, ca := c.RGBA()

	inner := func(ctx context.Context, tx database.Tx) error {
		err = r.queries.UpdateTag(ctx, tx, appdb.UpdateTagParams{
			ID:     tagID,
			Name:   name,
			ColorR: int64(cr),
			ColorG: int64(cg),
			ColorB: int64(cb),
			ColorA: int64(ca),
		})
		return errors.Wrap(err, "could not UpdateTag")
	}

	if tx == nil {
		err = errors.Wrap(database.TransactWithRetries(ctx, r.deps.Telemetry(), r.deps.Logger(), r.deps.DB(), &sql.TxOptions{}, inner), "could not TransactWithRetries")
	} else {
		err = inner(ctx, tx)
	}
	return err
}

func (r *AppSqliteRepository) DeleteTag(ctx context.Context, tagID int64, tx database.Tx) (err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "DeleteTag")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling DeleteTag", keys.TagID, tagID)

	inner := func(ctx context.Context, tx database.Tx) error {
		err = r.queries.DeleteTag(ctx, tx, tagID)
		return errors.Wrap(err, "could not DeleteTag")
	}

	if tx == nil {
		err = errors.Wrap(database.TransactWithRetries(ctx, r.deps.Telemetry(), r.deps.Logger(), r.deps.DB(), &sql.TxOptions{}, inner), "could not TransactWithRetries")
	} else {
		err = inner(ctx, tx)
	}
	return err
}

func (r *AppSqliteRepository) GetAllTags(ctx context.Context, tx database.Tx) (_ []*TagDBData, err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "GetAllTags")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling GetAllTags")

	tags, err := r.queries.GetAllTags(ctx, r.db(tx))
	if err != nil {
		return nil, errors.Wrap(err, "could not GetAllTags")
	}

	tagDBData := make([]*TagDBData, 0, len(tags))
	for _, t := range tags {
		skills, err := r.GetAllTagSkills(ctx, t.ID, tx)
		if err != nil && !errors.Is(err, database.ErrNoRows) {
			return nil, errors.Wrap(err, "could not GetAllTagSkills")
		}

		tagDBData = append(tagDBData, &TagDBData{
			Tag:    t,
			Skills: skills,
		})
	}

	return tagDBData, nil
}

func (r *AppSqliteRepository) GetAllTagSkills(ctx context.Context, tagID int64, tx database.Tx) (_ []TagSkill, err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "GetAllTagSkills")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling GetAllTagSkills")

	skills, err := r.queries.GetAllTagSkills(ctx, r.db(tx), tagID)
	if err != nil {
		return skills, err
	}

	return skills, nil
}

func (r *AppSqliteRepository) UpsertTagSkill(ctx context.Context, tagID, skillID, skillLevel int64, tx database.Tx) (skill TagSkill, err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "UpsertTagSkill")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling UpsertTagSkill", keys.TagID, tagID)

	inner := func(ctx context.Context, tx database.Tx) error {
		skill, err = r.queries.UpsertTagSkill(ctx, tx, appdb.UpsertTagSkillParams{
			TagID:      tagID,
			SkillID:    skillID,
			SkillLevel: skillLevel,
		})
		return errors.Wrap(err, "could not UpsertCharacterSkill")
	}

	if tx == nil {
		err = errors.Wrap(database.TransactWithRetries(ctx, r.deps.Telemetry(), r.deps.Logger(), r.deps.DB(), &sql.TxOptions{}, inner), "could not TransactWithRetries")
	} else {
		err = inner(ctx, tx)
	}
	return skill, err
}

func (r *AppSqliteRepository) DeleteTagSkills(ctx context.Context, tagID int64, skillIDs []int64, tx database.Tx) (err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.app", "DeleteTagSkills")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling DeleteTagSkills", keys.TagID, tagID, keys.SkillID, skillIDs)

	inner := func(ctx context.Context, tx database.Tx) error {
		err = r.queries.DeleteTagSkills(ctx, tx, appdb.DeleteTagSkillsParams{
			TagID:    tagID,
			SkillIds: skillIDs,
		})
		return errors.Wrap(err, "could not DeleteTagSkills")
	}

	if tx == nil {
		err = errors.Wrap(database.TransactWithRetries(ctx, r.deps.Telemetry(), r.deps.Logger(), r.deps.DB(), &sql.TxOptions{}, inner), "could not TransactWithRetries")
	} else {
		err = inner(ctx, tx)
	}
	return err
}
