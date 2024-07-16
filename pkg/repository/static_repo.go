package repository

import (
	"context"
	"strings"

	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"

	"github.com/kava-forge/eve-alts/pkg/database"
	"github.com/kava-forge/eve-alts/pkg/repository/internal/static"
	"github.com/kava-forge/eve-alts/pkg/telemetry"
)

type BatchGetSkillNamesRow = static.BatchGetSkillNamesRow

//counterfeiter:generate . StaticData
type StaticData interface {
	GetTableNames(ctx context.Context, tx database.Tx) ([]string, error)
	CleanTranslations(ctx context.Context, tx database.Tx) error
	GetSkillName(ctx context.Context, skillID int64, tx database.Tx) (string, error)
	GetSkillIDByName(ctx context.Context, skillName string, tx database.Tx) (int64, error)
	BatchGetSkillNames(ctx context.Context, skillIDs []int64, tx database.Tx) ([]BatchGetSkillNamesRow, error)
}

type staticDependencies interface {
	StaticDB() database.Connection
	Logger() logging.Logger
	Telemetry() *telemetry.Telemeter
}

type StaticSqliteRepository struct {
	deps    staticDependencies
	queries *static.Queries
}

var _ StaticData = (*StaticSqliteRepository)(nil)

func NewStaticData(deps staticDependencies) *StaticSqliteRepository {
	return &StaticSqliteRepository{
		deps:    deps,
		queries: static.New(),
	}
}

func (r *StaticSqliteRepository) db(tx database.Tx) static.DBTX {
	if tx != nil {
		return tx
	}
	return r.deps.StaticDB()
}

func (r *StaticSqliteRepository) GetTableNames(ctx context.Context, tx database.Tx) (_ []string, err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.static", "GetTableNames")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling GetTableNames")

	names, err := r.queries.GetTableNames(ctx, r.db(tx))
	if err != nil {
		return nil, err
	}

	return names, nil
}

func (r *StaticSqliteRepository) CleanTranslations(ctx context.Context, tx database.Tx) (err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.static", "CleanTranslations")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling CleanTranslations")

	return r.queries.CleanTranslations(ctx, r.db(tx))
}

func (r *StaticSqliteRepository) GetSkillName(ctx context.Context, skillID int64, tx database.Tx) (_ string, err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.static", "GetSkillName")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling GetSkillName")

	name, err := r.queries.GetSkillName(ctx, r.db(tx), static.GetSkillNameParams{
		SkillID:  skillID,
		Language: "en",
	})
	if err != nil {
		return "", err
	}

	return name, nil
}

func (r *StaticSqliteRepository) GetSkillIDByName(ctx context.Context, skillName string, tx database.Tx) (_ int64, err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.static", "GetSkillIDByName")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling GetSkillIDByName", "skill_name", skillName)

	id, err := r.queries.GetSkillIDFromName(ctx, r.db(tx), static.GetSkillIDFromNameParams{
		SkillNameLower: strings.ToLower(skillName),
		Language:       "en",
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *StaticSqliteRepository) BatchGetSkillNames(ctx context.Context, skillIDs []int64, tx database.Tx) (_ []BatchGetSkillNamesRow, err error) {
	ctx, span := telemetry.StartSpan(ctx, r.deps.Telemetry(), "repository.static", "BatchGetSkillNames")
	defer telemetry.EndSpan(span, &err)

	logger := r.deps.Logger()
	level.Debug(logger).Message("calling BatchGetSkillNames")

	rows, err := r.queries.BatchGetSkillNames(ctx, r.db(tx), static.BatchGetSkillNamesParams{
		SkillIds: skillIDs,
		Language: "en",
	})
	if err != nil {
		return nil, err
	}

	return rows, nil
}
