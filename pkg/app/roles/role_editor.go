package roles

import (
	"context"
	"database/sql"
	"fmt"
	"image/color"
	"math/rand"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kava-forge/eve-alts/lib/logging"

	"github.com/kava-forge/eve-alts/pkg/app/apperrors"
	"github.com/kava-forge/eve-alts/pkg/app/bindings"
	"github.com/kava-forge/eve-alts/pkg/app/colors"
	"github.com/kava-forge/eve-alts/pkg/database"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/operators"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

func populateRoleData(
	nameInp *widget.Entry,
	colorSwatch *colors.TappableColorSwatch,
	colorInp *dialog.ColorPickerDialog,
	labelInp *widget.Entry,
	operatorInp *widget.Select,
	tagSel *TagSelector,
	roleData bindings.DataProxy[*repository.RoleDBData],
) error {
	role, err := roleData.Get()
	if err != nil {
		return errors.Wrap(err, "unable to load role data")
	}

	nameInp.Text = role.Role.Name
	nameInp.Refresh()

	colorSwatch.SetColor(role.Color())
	colorSwatch.Refresh()
	colorInp.SetColor(role.Color())

	labelInp.Text = role.Role.Label
	labelInp.Refresh()

	operatorInp.Selected = role.Role.Operator
	operatorInp.Refresh()

	tagSel.Clear()
	for _, t := range role.Tags {
		tagSel.Select(t.ID)
	}

	return nil
}

func NewRoleEditor(deps dependencies, a fyne.App, title string, roles *bindings.DataList[*repository.RoleDBData], tags *bindings.DataList[*repository.TagDBData], roleData bindings.DataProxy[*repository.RoleDBData], onClose func()) fyne.Window {
	logger := logging.With(deps.Logger(), keys.Component, "RoleEditor")

	w := a.NewWindow(fmt.Sprintf("EVE Alts - %s", title))
	w.CenterOnScreen()
	w.Resize(fyne.Size{Width: 800, Height: 600})
	w.SetOnClosed(func() {
		if onClose != nil {
			onClose()
		}
	})

	nameInp := widget.NewEntry()
	labelInp := widget.NewEntry()
	operatorInp := widget.NewSelect(operators.OperatorNames(), func(string) {})

	colorSwatch := colors.NewTappableColorSwatch(deps.Logger(), color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255}) //nolint:gosec // not security related
	colorInp := dialog.NewColorPicker("Tag Color", "Pick a color for the role", func(c color.Color) {
		colorSwatch.SetColor(c)
	}, w)
	colorInp.Advanced = true
	colorSwatch.OnTapped = func(pe *fyne.PointEvent) { colorInp.Show() }

	tagSel := NewTagSelector(deps, w, tags)

	if roleData != nil {
		if err := populateRoleData(nameInp, colorSwatch, colorInp, labelInp, operatorInp, tagSel, roleData); err != nil {
			apperrors.Show(logger, w, apperrors.Error(
				"Could not load role data",
				apperrors.WithCause(err),
			), nil)
		}
	}

	form := widget.NewForm(
		widget.NewFormItem("Role Name", nameInp),
		widget.NewFormItem("Role Color", colorSwatch),
		widget.NewFormItem("Role Label", labelInp),
		widget.NewFormItem("Operator", operatorInp),
		widget.NewFormItem("Tags", tagSel.TagSet),
	)
	form.OnCancel = w.Close
	form.OnSubmit = func() {
		ctx := context.Background()

		tagIDs := tagSel.Selected()

		op, err := operators.ParseOperator(operatorInp.Selected)
		if err != nil {
			apperrors.Show(logger, w, apperrors.Error(
				"Unrecognized role operator",
				apperrors.WithCause(err),
			), nil)
		}

		if roleData == nil {
			var dbRole repository.Role
			var dbTags []repository.Tag
			if err := database.TransactWithRetries(ctx, deps.Telemetry(), logger, deps.DB(), &sql.TxOptions{}, func(ctx context.Context, tx database.Tx) error {
				var err error

				if dbTags == nil {
					dbTags = make([]repository.Tag, 0, 0) // TODO
				} else {
					dbTags = dbTags[:0]
				}

				dbRole, err = deps.AppRepo().InsertRole(ctx, nameInp.Text, labelInp.Text, op, colorSwatch.Color(), tx)
				if err != nil {
					return errors.Wrap(err, "could not InsertRole")
				}

				for _, tid := range tagIDs {
					_, err := deps.AppRepo().UpsertRoleTag(ctx, dbRole.ID, tid, tx)
					if err != nil && !errors.Is(err, database.ErrNoRows) {
						return errors.Wrap(err, "could not UpsertRoleTag", keys.TagID, tid)
					}
				}

				return nil
			}); err != nil {
				apperrors.Show(logger, w, apperrors.Error(
					"Could not create role",
					apperrors.WithCause(err),
				), nil)
				return
			}

			dbTags, err := deps.AppRepo().GetAllRoleTags(ctx, dbRole.ID, nil)
			if err != nil {
				apperrors.Show(logger, w, apperrors.Error(
					"Could not fetch role tags",
					apperrors.WithCause(err),
				), nil)
			}

			if err := roles.Append(&repository.RoleDBData{
				Role:     dbRole,
				Operator: operators.Operator(dbRole.Operator),
				Tags:     dbTags,
			}); err != nil {
				apperrors.Show(logger, w, apperrors.Error(
					"Could not append role",
					apperrors.WithCause(err),
				), nil)
			}
			w.Close()
		} else { // edit existing
			roleP, err := roleData.Get()
			if err != nil {
				apperrors.Show(logger, w, apperrors.Error(
					"Could not find role data",
					apperrors.WithCause(err),
				), nil)
				return
			}

			logger := logging.With(logger, keys.RoleID, roleP.Role.ID) //nolint:govet // intentional

			var dbRole repository.Role
			var dbTags []repository.Tag
			if err := database.TransactWithRetries(ctx, deps.Telemetry(), logger, deps.DB(), &sql.TxOptions{}, func(ctx context.Context, tx database.Tx) error {
				role := *roleP

				tagIDMap := make(map[int64]bool, len(tagIDs))
				for _, tid := range tagIDs {
					tagIDMap[tid] = true
				}
				toDelete := make([]int64, 0, 10)
				for _, t := range role.Tags {
					if !tagIDMap[t.ID] {
						toDelete = append(toDelete, t.ID)
					}
				}

				c := colorSwatch.Color()
				err = deps.AppRepo().UpdateRole(ctx, role.Role.ID, nameInp.Text, labelInp.Text, op, c, tx)
				if err != nil {
					return errors.Wrap(err, "could not UpdateRole")
				}
				role.Role.Name = nameInp.Text
				cr, cg, cb, ca := c.RGBA()
				role.Role.ColorR = int64(cr)
				role.Role.ColorG = int64(cg)
				role.Role.ColorB = int64(cb)
				role.Role.ColorA = int64(ca)

				for _, tid := range tagIDs {
					_, err := deps.AppRepo().UpsertRoleTag(ctx, role.Role.ID, tid, tx)
					if err != nil && !errors.Is(err, database.ErrNoRows) {
						return errors.Wrap(err, "could not UpsertRoleTag", keys.TagID, tid, keys.RoleID, role.Role.ID)
					}
				}

				if len(toDelete) > 0 {
					if err := deps.AppRepo().DeleteRoleTags(ctx, role.Role.ID, toDelete, tx); err != nil {
						return errors.Wrap(err, "could not DeleteRoleTags")
					}
				}

				dbRole = role.Role

				return nil
			}); err != nil {
				apperrors.Show(logger, w, apperrors.Error(
					"Could not save tag",
					apperrors.WithCause(err),
				), nil)
				return
			}

			dbTags, err = deps.AppRepo().GetAllRoleTags(ctx, dbRole.ID, nil)
			if err != nil {
				apperrors.Show(logger, w, apperrors.Error(
					"Could not fetch role tags",
					apperrors.WithCause(err),
				), nil)
			}

			if err := roleData.Set(&repository.RoleDBData{
				Role:     dbRole,
				Operator: operators.Operator(dbRole.Operator),
				Tags:     dbTags,
			}); err != nil {
				apperrors.Show(logger, w, apperrors.Error(
					"Could not set role",
					apperrors.WithCause(err),
				), nil)
			}
			w.Close()
		}
	}

	w.SetContent(form)

	return w
}
