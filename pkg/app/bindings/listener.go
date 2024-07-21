package bindings

import (
	"fyne.io/fyne/v2/data/binding"

	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/pkg/panics"
)

func NewListener(logger logging.Logger, cb func()) binding.DataListener {
	return binding.NewDataListener(func() {
		go func() {
			defer panics.Handler(logger)
			cb()
		}()
	})
}
