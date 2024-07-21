package bindings

import (
	"fyne.io/fyne/v2/data/binding"
)

func NewListener(cb func()) binding.DataListener {
	return binding.NewDataListener(func() {
		go cb()
	})
}
