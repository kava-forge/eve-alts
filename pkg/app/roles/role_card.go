package roles

import (
	"context"
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"

	"github.com/kava-forge/eve-alts/pkg/app/apperrors"
	"github.com/kava-forge/eve-alts/pkg/app/bindings"
	"github.com/kava-forge/eve-alts/pkg/app/colors"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

type RoleCard struct {
	widget.BaseWidget

	deps   dependencies
	role   bindings.DataProxy[*repository.RoleDBData]
	parent fyne.Window

	NameLabel    *widget.RichText
	ColorSwatch  *colors.ColorSwatch
	EditButton   *widget.Button
	DeleteButton *widget.Button

	update *sync.RWMutex
}

func NewRoleCard(deps dependencies, parent fyne.Window, dataRole bindings.DataProxy[*repository.RoleDBData], deleteFunc func(c *RoleCard), editFunc func(bindings.DataProxy[*repository.RoleDBData], func())) *RoleCard {
	logger := logging.With(deps.Logger(), keys.Component, "RoleCard")

	role, err := dataRole.Get()
	if err != nil {
		apperrors.Show(logger, parent, apperrors.Error(
			"Could not find role data",
			apperrors.WithCause(err),
		), nil)
		return nil
	}

	// logger = logging.With(logger, keys.TagID, tag.Tag.ID, keys.TagName, tag.Tag.Name)

	cc := &RoleCard{
		deps:   deps,
		role:   dataRole,
		parent: parent,

		NameLabel:   widget.NewRichTextWithText(role.Role.Name),
		ColorSwatch: colors.NewColorSwatch(deps.Logger(), role.Color()),
		// EditButton:   widget.NewButtonWithIcon("edit", theme.SettingsIcon(), nil),
		// DeleteButton: widget.NewButtonWithIcon("delete", theme.DeleteIcon(), nil),
		EditButton:   widget.NewButton("edit", nil),
		DeleteButton: widget.NewButton("delete", nil),

		update: &sync.RWMutex{},
	}
	cc.ExtendBaseWidget(cc)

	cc.refreshStyle()
	cc.ColorSwatch.SetCornerRadius(theme.InnerPadding() / 2)

	cc.EditButton.OnTapped = cc.editRole(editFunc)
	cc.DeleteButton.OnTapped = cc.deleteRole(deleteFunc)
	cc.DeleteButton.Importance = widget.DangerImportance

	dataRole.AddListener(bindings.NewListener(cc.redraw))

	return cc
}

func (c *RoleCard) Parent() fyne.Window {
	return c.parent
}

func (c *RoleCard) refreshStyle() {
	defer c.NameLabel.Refresh()

	c.NameLabel.Wrapping = fyne.TextWrapOff

	textColorName := theme.ColorNameForeground

	var darkText, lightText fyne.ThemeColorName
	if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantLight {
		darkText = theme.ColorNameForeground
		lightText = colors.ColorNameInvertedForeground
	} else {
		lightText = theme.ColorNameForeground
		darkText = colors.ColorNameInvertedForeground
	}

	if colors.UseDarkText(c.ColorSwatch.Color()) {
		textColorName = darkText
	} else {
		textColorName = lightText
	}

	for _, segi := range c.NameLabel.Segments {
		if seg, ok := segi.(*widget.TextSegment); ok {
			seg.Style = widget.RichTextStyle{
				Alignment: fyne.TextAlignLeading,
				ColorName: textColorName,
				Inline:    true,
				SizeName:  theme.SizeNameText,
				TextStyle: fyne.TextStyle{Bold: true},
			}
		}
		segi.Visual().Refresh()
	}
}

func (c *RoleCard) CreateRenderer() fyne.WidgetRenderer {
	return c
}

func (c *RoleCard) SetText(text string) {
	c.SetTextAt(0, text)
}

func (c *RoleCard) SetTextAt(idx int, text string) {
	c.NameLabel.Segments[idx].(*widget.TextSegment).Text = text
	c.NameLabel.Refresh()
}

func (c *RoleCard) redraw() {
	c.update.Lock()
	defer c.Refresh()
	defer c.update.Unlock()

	logger := logging.With(c.deps.Logger(), keys.Component, "RoleCard.redraw")

	role, err := c.role.Get()
	if err != nil || role == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find role data",
			apperrors.WithCause(err),
		), nil)
		return
	}

	logger = logging.With(logger, keys.RoleID, role.Role.ID, keys.RoleName, role.Role.Name)

	level.Debug(logger).Message("refreshing role card", "color", fmt.Sprintf("%#v", role.Color()))

	c.SetText(role.Role.Name)
	c.ColorSwatch.SetColor(role.Color())
	c.refreshStyle()
}

func (c *RoleCard) RoleID() int64 {
	c.update.RLock()
	defer c.update.RUnlock()
	logger := logging.With(c.deps.Logger(), keys.Component, "RoleCard.RoleID")

	role, err := c.role.Get()
	if err != nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find role data",
			apperrors.WithCause(err),
		), nil)
		return 0
	}
	return role.Role.ID
}

func (c *RoleCard) RefreshDataWith(role *repository.RoleDBData) error {
	return c.role.Set(role)
}

func (c *RoleCard) editRole(editFunc func(bindings.DataProxy[*repository.RoleDBData], func())) func() {
	return func() {
		c.EditButton.Disable()

		editFunc(c.role, func() {
			c.EditButton.Enable()
		})
	}
}

func (c *RoleCard) deleteRole(callback func(*RoleCard)) func() {
	logger := logging.With(c.deps.Logger(), keys.Component, "RoleCard.deleteRole")
	return func() {
		ctx := context.Background()

		c.DeleteButton.Disable()

		role, err := c.role.Get()
		if err != nil {
			apperrors.Show(logger, c.parent, apperrors.Error(
				"Could not find role data",
				apperrors.WithCause(err),
			), nil)
			return
		}

		logger := logging.With(logger, keys.RoleID, role.Role.ID, keys.RoleName, role.Role.Name) //nolint:govet // intentional

		conf := dialog.NewConfirm("Delete Role?", fmt.Sprintf("Are you sure you want to delete the role '%s'?", role.Role.Name), func(ok bool) {
			defer c.DeleteButton.Enable()
			if !ok {
				return
			}

			if err := c.deps.AppRepo().DeleteRole(ctx, c.RoleID(), nil); err != nil {
				apperrors.Show(logger, c.parent, apperrors.Error(
					"Could not delete role",
					apperrors.WithCause(err),
				), nil)
			} else {
				callback(c)
			}
		}, c.parent)
		conf.SetConfirmImportance(widget.DangerImportance)
		conf.Show()
	}
}

// The WidgetRenderer interface

func (c *RoleCard) Destroy() {}

func (c *RoleCard) Layout(sz fyne.Size) {
	c.update.RLock()
	defer c.update.RUnlock()

	fontSize := fyne.CurrentApp().Settings().Theme().Size(theme.SizeNameText)

	c.ColorSwatch.Move(fyne.Position{X: 0, Y: 0})
	c.ColorSwatch.Resize(sz)

	c.NameLabel.Move(fyne.Position{X: 6, Y: 0})
	c.NameLabel.Refresh()

	editLabelSz := fyne.MeasureText(c.EditButton.Text, fontSize, fyne.TextStyle{})
	c.EditButton.Resize(fyne.Size{Width: editLabelSz.Width + 2*theme.InnerPadding(), Height: editLabelSz.Height + theme.InnerPadding()})
	c.EditButton.Alignment = widget.ButtonAlignCenter
	rbsz := c.EditButton.Size()
	c.EditButton.Move(fyne.Position{X: sz.Width - rbsz.Width - theme.Padding(), Y: sz.Height - rbsz.Height - theme.Padding()})

	deleteLabelSz := fyne.MeasureText(c.DeleteButton.Text, fontSize, fyne.TextStyle{})
	c.DeleteButton.Resize(fyne.Size{Width: deleteLabelSz.Width + 2*theme.InnerPadding(), Height: deleteLabelSz.Height + theme.InnerPadding()})
	c.DeleteButton.Alignment = widget.ButtonAlignCenter
	dbsz := c.DeleteButton.Size()
	c.DeleteButton.Move(fyne.Position{X: sz.Width - rbsz.Width - theme.Padding() - dbsz.Width - theme.Padding(), Y: sz.Height - dbsz.Height - theme.Padding()})
}

func (c *RoleCard) MinSize() fyne.Size {
	return fyne.Size{
		Height: 50,
		Width:  100,
	}
}

func (c *RoleCard) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		c.ColorSwatch,
		c.NameLabel,
		c.EditButton,
		c.DeleteButton,
	}
}

func (c *RoleCard) Refresh() {
	for _, o := range c.Objects() {
		o.Refresh()
	}
}
