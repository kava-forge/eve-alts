package telemetry

import (
	"github.com/kava-forge/eve-alts/lib/telemetry"
)

type Stats struct {
	RequestCount telemetry.Int64Counter
}

func NewStats(appName string, telemeter *Telemeter) (*Stats, error) {
	meter := telemeter.Meter(appName)

	var err error
	stats := &Stats{}

	stats.RequestCount, err = meter.Int64Counter("request_ct")
	if err != nil {
		return stats, err
	}

	return stats, nil
}
