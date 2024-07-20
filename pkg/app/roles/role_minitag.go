package roles

import (
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

type TagMiniTag struct {
	*minitag.MiniTag

	deps     dependencies
	parent   fyne.Window
	tag      bindings.DataProxy[*repository.TagDBData]
	selected bindings.DataProxy[bool]
}

func NewTagMiniTag(deps dependencies, parent fyne.Window, tagData bindings.DataProxy[*repository.TagDBData], selected bindings.DataProxy[bool]) *TagMiniTag {
	logger := logging.With(deps.Logger(), keys.Component, "TagMiniTag")

	tag, err := tagData.Get()
	if err != nil {
		apperrors.Show(logger, parent, apperrors.Error(
			"Could not find tag data",
			apperrors.WithCause(err),
		), nil)
		return nil
	}

	tmt := &TagMiniTag{
		MiniTag:  minitag.New(logger, tag.Tag.Name, tag.Color(), theme.SizeNameText),
		deps:     deps,
		parent:   parent,
		tag:      tagData,
		selected: selected,
	}

	tmt.MiniTag.UnDimmedBold = true

	tmt.redraw()

	tagData.AddListener(binding.NewDataListener(tmt.redraw))
	selected.AddListener(binding.NewDataListener(tmt.redraw))

	return tmt
}

func (c *TagMiniTag) redraw() {
	defer c.Refresh()

	logger := logging.With(c.deps.Logger(), keys.Component, "TagMiniTag.redraw")

	tag, err := c.tag.Get()
	if err != nil || tag == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find tag data",
			apperrors.WithCause(err),
		), nil)
		return
	}

	logger = logging.With(logger, keys.TagID, tag.Tag.ID, keys.TagName, tag.Tag.Name)

	if tag.Tag.ID == 0 { // deleted
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

	c.SetText(tag.Tag.Name)
	c.ColorSwatch.SetColor(tag.Color())
	c.MiniTag.Dimmed = !selected
	c.RefreshStyle()
	c.MiniTag.Resize(c.MiniTag.Size())
}

func (c *TagMiniTag) ShouldShow() bool {
	logger := logging.With(c.deps.Logger(), keys.Component, "TagMiniTag.ShouldShow")

	tag, err := c.tag.Get()
	if err != nil || tag == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find tag data",
			apperrors.WithCause(err),
		), nil)
		return false
	}

	// logger = logging.With(logger, keys.TagID, tag.Tag.ID, keys.TagName, tag.Tag.Name)

	if tag.Tag.ID == 0 { // deleted
		return false
	}

	return true
}

func (c *TagMiniTag) SortKey() string {
	logger := logging.With(c.deps.Logger(), keys.Component, "TagMiniTag.SortKey")

	tag, err := c.tag.Get()
	if err != nil || tag == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find tag data",
			apperrors.WithCause(err),
		), nil)
		return ""
	}

	return tag.Tag.Name
}

var _ fyne.Tappable = (*TagMiniTag)(nil)

func (c *TagMiniTag) Tapped(_ *fyne.PointEvent) {
	logger := logging.With(c.deps.Logger(), keys.Component, "TagMiniTag.Tapped")
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
