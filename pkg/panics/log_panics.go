package panics

import (
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"
)

func Handler(logger logging.Logger) {
	if err := recover(); err != nil {
		level.Error(logger).Message("panic!", "error", err)
		panic(err)
	}
}
