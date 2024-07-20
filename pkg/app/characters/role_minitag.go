package characters

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"

	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"

	"github.com/kava-forge/eve-alts/pkg/app/apperrors"
	"github.com/kava-forge/eve-alts/pkg/app/bindings"
	"github.com/kava-forge/eve-alts/pkg/app/minitag"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

type RoleMiniTag struct {
	*minitag.MiniTag

	deps      dependencies
	parent    fyne.Window
	char      bindings.DataProxy[*repository.CharacterDBData]
	role      bindings.DataProxy[*repository.RoleDBData]
	tags      *bindings.DataList[*repository.TagDBData]
	knownTags map[int64]bool
	isMatch   bool
	missing   []string
}

func NewRoleMiniTag(deps dependencies, parent fyne.Window, char bindings.DataProxy[*repository.CharacterDBData], role bindings.DataProxy[*repository.RoleDBData], tags *bindings.DataList[*repository.TagDBData]) *RoleMiniTag {
	logger := logging.With(deps.Logger(), keys.Component, "RoleMiniTag")

	roleData, err := role.Get()
	if err != nil {
		apperrors.Show(logger, parent, apperrors.Error(
			"Could not find role data",
			apperrors.WithCause(err),
		), nil)
		return nil
	}

	cmt := &RoleMiniTag{
		MiniTag:   minitag.New(logger, roleData.Role.Name, roleData.Color(), theme.SizeNameCaptionText),
		deps:      deps,
		parent:    parent,
		char:      char,
		role:      role,
		tags:      tags,
		knownTags: map[int64]bool{},
	}

	cmt.addTag()
	cmt.redraw()

	char.AddListener(binding.NewDataListener(cmt.redraw))
	role.AddListener(binding.NewDataListener(cmt.redraw))
	tags.AddListener(binding.NewDataListener(func() {
		cmt.addTag()
		cmt.redraw()
	}))

	return cmt
}

func (c *RoleMiniTag) addTag() {
	logger := logging.With(c.deps.Logger(), keys.Component, "RoleMiniTag.addTag")

	for i := range c.tags.Length() {
		t, err := c.tags.GetValue(i)
		if err != nil {
			level.Error(logger).Err("could not retrieve tag child data", err)
			continue
		}

		if c.knownTags[t.Tag.ID] {
			continue
		}

		c.tags.Child(i).AddListener(binding.NewDataListener(c.redraw))

		c.knownTags[t.Tag.ID] = true
	}
}

func (c *RoleMiniTag) redraw() {
	defer c.Refresh()

	logger := logging.With(c.deps.Logger(), keys.Component, "RoleMiniTag.redraw")

	role, err := c.role.Get()
	if err != nil || role == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find role data",
			apperrors.WithCause(err),
		), nil)
		return
	}

	logger = logging.With(logger, keys.RoleID, role.Role.ID, keys.RoleName, role.Role.Name)

	char, err := c.char.Get()
	if err != nil || char == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find character data",
			apperrors.WithCause(err),
		), nil)
		return
	}

	logger = logging.With(logger, keys.CharacterID, char.Character.ID, keys.CharacterName, char.Character.Name)

	if role.Role.ID == 0 || char.Character.ID == 0 { // deleted
		return
	}

	tags, err := c.tags.Get()
	if err != nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find tag list data",
			apperrors.WithCause(err),
		), nil)
		return
	}

	isMatch, missing := CharacterMatchesRole(char, role, tags)
	level.Debug(logger).Message("role match?", "match", isMatch, "missing", missing)

	c.SetText(role.Role.Label)
	c.ColorSwatch.SetColor(role.Color())
	c.isMatch = isMatch
	c.MiniTag.Dimmed = !c.isMatch
	c.RefreshStyle()

	c.missing = c.missing[:0]
	for _, m := range missing {
		c.missing = append(c.missing, m.Name)
	}
}

func (c *RoleMiniTag) ShouldShow() bool {
	logger := logging.With(c.deps.Logger(), keys.Component, "RoleMiniTag.ShouldShow")

	role, err := c.role.Get()
	if err != nil || role == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find role data",
			apperrors.WithCause(err),
		), nil)
		return false
	}

	logger = logging.With(logger, keys.RoleID, role.Role.ID, keys.RoleName, role.Role.Name)

	char, err := c.char.Get()
	if err != nil || char == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find character data",
			apperrors.WithCause(err),
		), nil)
		return false
	}

	if role.Role.ID == 0 || char.Character.ID == 0 { // deleted
		return false
	}

	return c.isMatch
}

func (c *RoleMiniTag) SortKey() string {
	logger := logging.With(c.deps.Logger(), keys.Component, "RoleMiniTag.SortKey")

	role, err := c.role.Get()
	if err != nil || role == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find role data",
			apperrors.WithCause(err),
		), nil)
		return ""
	}

	if c.isMatch {
		return fmt.Sprintf(" %s", role.Role.Label)
	}
	return role.Role.Label
}
