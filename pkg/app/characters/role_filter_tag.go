package characters

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"

	"github.com/kava-forge/eve-alts/pkg/app/apperrors"
	"github.com/kava-forge/eve-alts/pkg/app/bindings"
	"github.com/kava-forge/eve-alts/pkg/app/minitag"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

type RoleFilterTag struct {
	*minitag.MiniTag

	deps     dependencies
	parent   fyne.Window
	role     bindings.DataProxy[*repository.RoleDBData]
	selected bindings.DataProxy[bool]

	update *sync.RWMutex
}

func NewRoleFilterTag(deps dependencies, parent fyne.Window, roleData bindings.DataProxy[*repository.RoleDBData], selected bindings.DataProxy[bool]) *RoleFilterTag {
	logger := logging.With(deps.Logger(), keys.Component, "RoleFilterTag")

	role, err := roleData.Get()
	if err != nil {
		apperrors.Show(logger, parent, apperrors.Error(
			"Could not find role data",
			apperrors.WithCause(err),
		), nil)
		return nil
	}

	tmt := &RoleFilterTag{
		MiniTag:  minitag.New(logger, role.Role.Name, role.Color(), theme.SizeNameText),
		deps:     deps,
		parent:   parent,
		role:     roleData,
		selected: selected,

		update: &sync.RWMutex{},
	}

	tmt.MiniTag.UnDimmedBold = true

	tmt.redraw()

	roleData.AddListener(bindings.NewListener(logger, tmt.redraw))
	selected.AddListener(bindings.NewListener(logger, tmt.redraw))

	return tmt
}

func (c *RoleFilterTag) redraw() {
	c.update.Lock()
	defer c.Refresh()
	defer c.update.Unlock()

	logger := logging.With(c.deps.Logger(), keys.Component, "RoleFilterTag.redraw")

	role, err := c.role.Get()
	if err != nil || role == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find role data",
			apperrors.WithCause(err),
		), nil)
		return
	}

	logger = logging.With(logger, keys.RoleID, role.Role.ID, keys.RoleName, role.Role.Name)

	if role.Role.ID == 0 { // deleted
		if !c.Hidden {
			c.Hide()
		}
	}

	selected, err := c.selected.Get()
	if err != nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find selected state data",
			apperrors.WithCause(err),
		), nil)
		return
	}

	level.Debug(logger).Message("new state", "new", selected)

	c.SetText(role.Role.Name)
	c.ColorSwatch.SetColor(role.Color())
	c.MiniTag.Dimmed = !selected
	c.MiniTag.RefreshStyle()
	c.MiniTag.Resize(c.MiniTag.Size())
}

func (c *RoleFilterTag) ShouldShow() bool {
	logger := logging.With(c.deps.Logger(), keys.Component, "RoleFilterTag.ShouldShow")

	role, err := c.role.Get()
	if err != nil || role == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find role data",
			apperrors.WithCause(err),
		), nil)
		return false
	}

	// logger = logging.With(logger, keys.TagID, tag.Tag.ID, keys.TagName, tag.Tag.Name)

	if role.Role.ID == 0 { // deleted
		return false
	}

	return true
}

func (c *RoleFilterTag) SortKey() string {
	logger := logging.With(c.deps.Logger(), keys.Component, "RoleFilterTag.SortKey")

	role, err := c.role.Get()
	if err != nil || role == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find role data",
			apperrors.WithCause(err),
		), nil)
		return ""
	}

	return role.Role.Name
}

var _ fyne.Tappable = (*RoleFilterTag)(nil)

func (c *RoleFilterTag) Tapped(_ *fyne.PointEvent) {
	logger := logging.With(c.deps.Logger(), keys.Component, "RoleFilterTag.Tapped")
	level.Debug(logger).Message("minitag tap")

	selected, err := c.selected.Get()
	if err != nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find selected state data",
			apperrors.WithCause(err),
		), nil)
		return
	}

	if err := c.selected.Set(!selected); err != nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not toggle tag selection",
			apperrors.WithCause(err),
		), nil)
	}
}
