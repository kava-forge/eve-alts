package characters

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"

	"github.com/kava-forge/eve-alts/pkg/app/apperrors"
	"github.com/kava-forge/eve-alts/pkg/app/bindings"
	"github.com/kava-forge/eve-alts/pkg/app/minitag"
	"github.com/kava-forge/eve-alts/pkg/keys"
	"github.com/kava-forge/eve-alts/pkg/repository"
)

type CharacterMiniTag struct {
	*minitag.MiniTag

	deps    dependencies
	parent  fyne.Window
	char    bindings.DataProxy[*repository.CharacterDBData]
	tag     bindings.DataProxy[*repository.TagDBData]
	isMatch bool
	missing []string

	update *sync.RWMutex
}

func NewCharacterMiniTag(deps dependencies, parent fyne.Window, char bindings.DataProxy[*repository.CharacterDBData], tag bindings.DataProxy[*repository.TagDBData]) *CharacterMiniTag {
	logger := logging.With(deps.Logger(), keys.Component, "CharacterMiniTag")

	tagData, err := tag.Get()
	if err != nil {
		apperrors.Show(logger, parent, apperrors.Error(
			"Could not find tag data",
			apperrors.WithCause(err),
		), nil)
		return nil
	}

	cmt := &CharacterMiniTag{
		MiniTag: minitag.New(logger, tagData.Tag.Name, tagData.Color(), theme.SizeNameCaptionText),
		deps:    deps,
		parent:  parent,
		char:    char,
		tag:     tag,

		update: &sync.RWMutex{},
	}

	cmt.redraw()

	char.AddListener(bindings.NewListener(cmt.redraw))
	tag.AddListener(bindings.NewListener(cmt.redraw))

	return cmt
}

func (c *CharacterMiniTag) redraw() {
	c.update.Lock()
	defer c.Refresh()
	defer c.update.Unlock()

	logger := logging.With(c.deps.Logger(), keys.Component, "CharacterMiniTag.redraw")

	tag, err := c.tag.Get()
	if err != nil || tag == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find tag data",
			apperrors.WithCause(err),
		), nil)
		return
	}

	logger = logging.With(logger, keys.TagID, tag.Tag.ID, keys.TagName, tag.Tag.Name)

	char, err := c.char.Get()
	if err != nil || char == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find character data",
			apperrors.WithCause(err),
		), nil)
		return
	}

	logger = logging.With(logger, keys.CharacterID, char.Character.ID, keys.CharacterName, char.Character.Name)

	if tag.Tag.ID == 0 || char.Character.ID == 0 { // deleted
		return
	}

	isMatch, missing := CharacterMatchesTag(char, tag)
	level.Debug(logger).Message("tag match?", "match", isMatch, "missing", missing)

	c.SetText(tag.Tag.Name)
	c.ColorSwatch.SetColor(tag.Color())
	c.isMatch = isMatch
	c.MiniTag.Dimmed = !c.isMatch
	c.MiniTag.RefreshStyle()

	ids := make([]int64, 0, len(missing))
	for _, sk := range missing {
		ids = append(ids, sk.SkillID)
	}

	names, err := c.deps.StaticRepo().BatchGetSkillNames(context.Background(), ids, nil)
	if err != nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find load skill names",
			apperrors.WithCause(err),
		), nil)
		return
	}

	nameMap := make(map[int64]string, len(names))
	for _, name := range names {
		nameMap[name.SkillID] = name.SkillName
	}

	c.missing = c.missing[:0]
	for _, sk := range missing {
		c.missing = append(c.missing, fmt.Sprintf("%s %d", nameMap[sk.SkillID], sk.SkillLevel))
	}
}

func (c *CharacterMiniTag) ShouldShow() bool {
	c.update.RLock()
	defer c.update.RUnlock()

	logger := logging.With(c.deps.Logger(), keys.Component, "CharacterMiniTag.ShouldShow")

	tag, err := c.tag.Get()
	if err != nil || tag == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find tag data",
			apperrors.WithCause(err),
		), nil)
		return false
	}

	logger = logging.With(logger, keys.TagID, tag.Tag.ID, keys.TagName, tag.Tag.Name)

	char, err := c.char.Get()
	if err != nil || char == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find character data",
			apperrors.WithCause(err),
		), nil)
		return false
	}

	if tag.Tag.ID == 0 || char.Character.ID == 0 { // deleted
		return false
	}

	return true
}

func (c *CharacterMiniTag) SortKey() string {
	c.update.RLock()
	defer c.update.RUnlock()

	logger := logging.With(c.deps.Logger(), keys.Component, "CharacterMiniTag.SortKey")

	tag, err := c.tag.Get()
	if err != nil || tag == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find tag data",
			apperrors.WithCause(err),
		), nil)
		return ""
	}

	if c.isMatch {
		return fmt.Sprintf(" %s", tag.Tag.Name)
	}
	return tag.Tag.Name
}

var _ fyne.Tappable = (*CharacterMiniTag)(nil)

func (c *CharacterMiniTag) Tapped(_ *fyne.PointEvent) {
	logger := logging.With(c.deps.Logger(), keys.Component, "CharacterMiniTag.Tapped")
	level.Debug(logger).Message("minitag tap")
	if !c.isMatch {
		list := widget.NewRichTextWithText(strings.Join(c.missing, "\n"))
		// scr := container.NewVScroll(list)
		estTime := widget.NewLabel("Estimated train time: TBD")
		copyButton := widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), func() {
			c.parent.Clipboard().SetContent(list.String())
		})
		data := container.New(layout.NewVBoxLayout(), list, layout.NewSpacer(), container.New(layout.NewHBoxLayout(), estTime, layout.NewSpacer(), copyButton))
		d := dialog.NewCustom(fmt.Sprintf("Missing Skills for %s", c.NameLabel.String()), "Close", data, c.parent)
		d.Show()
	}
}
