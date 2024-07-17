package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/kava-forge/eve-alts/lib/deferutil"
	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/logging/level"
	"github.com/kava-forge/eve-alts/pkg/database"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

// build time variables
var (
	AppName      string
	BuildDate    string
	BuildVersion string
	BuildSHA     string
)

func main() {
	err := run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", AppName, err)
		os.Exit(1)
	}
}

func run() error {
	var databaseFile string

	fs := flag.NewFlagSet(AppName, flag.ContinueOnError)
	fs.StringVar(&databaseFile, "database", "", "The database to clean")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return errors.Wrap(err, "could not parse args")
	}

	ctx := context.Background()

	deps, err := CreateDependencies(ctx, AppName, BuildVersion, databaseFile)
	if err != nil {
		return errors.Wrap(err, "could not create dependencies")
	}
	defer deferutil.CheckDefer(func() error { return deps.Telemetry().Shutdown(ctx) })

	level.Debug(deps.Logger()).Message("database file", "fname", databaseFile, "cwd", os.Getenv("PWD"), "args", os.Args)

	level.Debug(deps.Logger()).Message("creating repo")
	repo := repository.NewStaticData(deps)
	level.Debug(deps.Logger()).Message("getting table names")
	names, err := repo.GetTableNames(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "could not get table names")
	}
	level.Debug(deps.Logger()).Message("got tables", "names", names)

	for _, n := range names {
		switch n {
		case "trnTranslations":
			continue
		default:
			level.Debug(deps.Logger()).Message("dropping table", "table_name", n)
			if _, err := deps.StaticDB().ExecContext(ctx, fmt.Sprintf("DROP TABLE %s", n)); err != nil {
				return errors.Wrap(err, "could not drop table", "table_name", n)
			}
		}
	}

	if _, err := deps.StaticDB().ExecContext(ctx, `
CREATE INDEX IF NOT EXISTS "idx_translations_by_name" 
ON trnTranslations ("tcID", "languageID", LOWER("text"))
WHERE "tcID" = 8;
		`); err != nil {
		return errors.Wrap(err, "could not create index")
	}

	if err := repo.CleanTranslations(ctx, nil); err != nil {
		return errors.Wrap(err, "could not CleanTranslations")
	}

	level.Debug(deps.Logger()).Message("vacuuming database")
	if _, err := deps.StaticDB().ExecContext(ctx, "VACUUM"); err != nil {
		return errors.Wrap(err, "could not vacuum")
	}

	return deps.staticDB.(*database.WrappedConnection).DB.Close()
}
