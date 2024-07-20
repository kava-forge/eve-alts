package roles

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/kava-forge/eve-alts/pkg/app/bindings"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

func NewAddRoleButton(deps dependencies, roles *bindings.DataList[*repository.RoleDBData], tags *bindings.DataList[*repository.TagDBData]) *widget.Button {
	button := widget.NewButtonWithIcon("Add Role", theme.ContentAddIcon(), func() {
		a := fyne.CurrentApp()

		w := NewRoleEditor(deps, a, "Add Role", roles, tags, nil, nil)
		w.Show()
	})

	return button
}
