package characters

import (
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/pkg/database"
	"github.com/kava-forge/eve-alts/pkg/esi"
	"github.com/kava-forge/eve-alts/pkg/repository"
	"github.com/kava-forge/eve-alts/pkg/telemetry"
)

type dependencies interface {
	DB() database.Connection
	StaticDB() database.Connection
	Logger() logging.Logger
	ESIClient() esi.Client

	Telemetry() *telemetry.Telemeter
	Stats() *telemetry.Stats

	AppRepo() repository.AppData
	StaticRepo() repository.StaticData
}
