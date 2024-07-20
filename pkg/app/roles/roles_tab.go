package roles

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"

	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"

	"github.com/kava-forge/eve-alts/pkg/app/apperrors"
	"github.com/kava-forge/eve-alts/pkg/app/bindings"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

func NewRolesTab(deps dependencies, parent fyne.Window, tagsData *bindings.DataList[*repository.TagDBData]) (fyne.CanvasObject, *bindings.DataList[*repository.RoleDBData]) {
	logger := logging.With(deps.Logger(), keys.Component, "MainWindow.RolesTab")

	lout := layout.NewGridWrapLayout(fyne.Size{Width: 250, Height: 50})
	ctn := container.New(lout)

	rolesData := bindings.NewDataList[*repository.RoleDBData]()
	rolesData.AddListener(binding.NewDataListener(func() {
		defer func() {
			level.Debug(logger).Message("roles listener done")
		}()
		// defer ctn.Refresh()

		level.Debug(logger).Message("roles listener start")

		level.Debug(logger).Message("refreshing roles shown")
		rolesShown := make(map[int64]int, rolesData.Length())
		toDelete := make(map[int64]int, rolesData.Length())

		for i, o := range ctn.Objects {
			if c, ok := o.(*RoleCard); ok {
				rolesShown[c.RoleID()] = i
				toDelete[c.RoleID()] = i
			}
		}

		rolesList, err := rolesData.Get()
		if err != nil {
			apperrors.Show(logger, parent, apperrors.Error(
				"Could not find role list data",
				apperrors.WithCause(err),
			), nil)
			return
		}

		level.Debug(logger).Message("tag list", "shown", rolesShown, "toDelete", toDelete)

		for i, role := range rolesList {
			if role == nil || role.Role.ID == 0 {
				continue
			}

			logger := logging.With(logger, keys.RoleID, role.Role.ID, keys.RoleName, role.Role.Name) //nolint:govet // intentional

			if idx := rolesShown[role.Role.ID]; idx != 0 {
				level.Debug(logger).Message("refreshing tag", keys.RoleID, role.Role.ID, "index", idx)
				delete(toDelete, role.Role.ID)
				if err := ctn.Objects[idx].(*RoleCard).RefreshDataWith(role); err != nil {
					apperrors.Show(logger, parent, apperrors.Error(
						"Could not set roles data",
						apperrors.WithCause(err),
					), nil)
				}
				continue
			}

			cc := NewRoleCard(deps, parent, rolesData.Child(i), func(c *RoleCard) {
				level.Debug(logger).Message("removing role card")
				role, err := rolesData.GetValue(i) //nolint:govet // intentional
				if err != nil {
					apperrors.Show(logger, c.Parent(), apperrors.Error(
						"Could not find role data",
						apperrors.WithCause(err),
					), nil)
					return
				}
				role.Role.ID = 0
				if err := rolesData.SetValue(i, role); err != nil {
					apperrors.Show(logger, c.Parent(), apperrors.Error(
						"Could not delete role",
						apperrors.WithCause(err),
					), nil)
				}
				ctn.Remove(c)
			}, func(roleData bindings.DataProxy[*repository.RoleDBData], onClose func()) {
				w := NewRoleEditor(deps, fyne.CurrentApp(), "Edit Role", rolesData, tagsData, roleData, onClose)
				w.Show()
			})
			level.Debug(logger).Message("adding new role", keys.RoleID, role.Role.ID, "cc", fmt.Sprintf("%#v", cc))
			ctn.Add(cc)
		}

		for roleID, idx := range toDelete {
			level.Debug(logger).Message("removing widget", keys.RoleID, roleID, "index", idx)
			ctn.Remove(ctn.Objects[idx])
		}
	}))

	ctn.Add(NewAddRoleButton(deps, rolesData, tagsData))

	return container.NewVScroll(ctn), rolesData
}
