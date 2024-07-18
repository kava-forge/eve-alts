package tags

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/kava-forge/eve-alts/pkg/app/bindings"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

func NewAddTagButton(deps dependencies, tags *bindings.DataList[*repository.TagDBData]) *widget.Button {
	button := widget.NewButtonWithIcon("Add Tag", theme.ContentAddIcon(), func() {
		a := fyne.CurrentApp()

		w := NewTagEditor(deps, a, "Add Tag", tags, nil, nil)
		w.Show()
	})

	return button
}
