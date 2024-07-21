package tags

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

type TagCard struct {
	widget.BaseWidget

	deps   dependencies
	tag    bindings.DataProxy[*repository.TagDBData]
	parent fyne.Window

	NameLabel    *widget.RichText
	ColorSwatch  *colors.ColorSwatch
	EditButton   *widget.Button
	DeleteButton *widget.Button

	update *sync.RWMutex
}

func NewTagCard(deps dependencies, parent fyne.Window, dataTag bindings.DataProxy[*repository.TagDBData], deleteFunc func(c *TagCard), editFunc func(bindings.DataProxy[*repository.TagDBData], func())) *TagCard {
	logger := logging.With(deps.Logger(), keys.Component, "TagCard")

	tag, err := dataTag.Get()
	if err != nil {
		apperrors.Show(logger, parent, apperrors.Error(
			"Could not find tag data",
			apperrors.WithCause(err),
		), nil)
		return nil
	}

	// logger = logging.With(logger, keys.TagID, tag.Tag.ID, keys.TagName, tag.Tag.Name)

	cc := &TagCard{
		deps:   deps,
		tag:    dataTag,
		parent: parent,

		NameLabel:   widget.NewRichTextWithText(tag.Tag.Name),
		ColorSwatch: colors.NewColorSwatch(deps.Logger(), tag.Color()),
		// EditButton:   widget.NewButtonWithIcon("edit", theme.SettingsIcon(), nil),
		// DeleteButton: widget.NewButtonWithIcon("delete", theme.DeleteIcon(), nil),
		EditButton:   widget.NewButton("edit", nil),
		DeleteButton: widget.NewButton("delete", nil),

		update: &sync.RWMutex{},
	}
	cc.ExtendBaseWidget(cc)

	cc.refreshStyle()
	cc.ColorSwatch.SetCornerRadius(theme.InnerPadding() / 2)

	cc.EditButton.OnTapped = cc.editTag(editFunc)
	cc.DeleteButton.OnTapped = cc.deleteTag(deleteFunc)
	cc.DeleteButton.Importance = widget.DangerImportance

	dataTag.AddListener(bindings.NewListener(logger, cc.redraw))

	return cc
}

func (c *TagCard) Parent() fyne.Window {
	return c.parent
}

func (c *TagCard) refreshStyle() {
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

func (c *TagCard) CreateRenderer() fyne.WidgetRenderer {
	return c
}

func (c *TagCard) setText(text string) {
	c.setTextAt(0, text)
}

func (c *TagCard) setTextAt(idx int, text string) {
	c.NameLabel.Segments[idx].(*widget.TextSegment).Text = text
	c.NameLabel.Refresh()
}

func (c *TagCard) redraw() {
	c.update.Lock()
	defer c.Refresh()
	defer c.update.Unlock()

	logger := logging.With(c.deps.Logger(), keys.Component, "TagCard.redraw")

	tag, err := c.tag.Get()
	if err != nil || tag == nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find tag data",
			apperrors.WithCause(err),
		), nil)
		return
	}

	logger = logging.With(logger, keys.TagID, tag.Tag.ID, keys.TagName, tag.Tag.Name)

	level.Debug(logger).Message("refreshing tag card", "color", fmt.Sprintf("%#v", tag.Color()))

	c.setText(tag.Tag.Name)
	c.ColorSwatch.SetColor(tag.Color())
	c.refreshStyle()
}

func (c *TagCard) TagID() int64 {
	c.update.RLock()
	defer c.update.RUnlock()

	logger := logging.With(c.deps.Logger(), keys.Component, "TagCard.TagID")

	tag, err := c.tag.Get()
	if err != nil {
		apperrors.Show(logger, c.parent, apperrors.Error(
			"Could not find tag data",
			apperrors.WithCause(err),
		), nil)
		return 0
	}
	return tag.Tag.ID
}

func (c *TagCard) RefreshDataWith(tag *repository.TagDBData) error {
	return c.tag.Set(tag)
}

func (c *TagCard) editTag(editFunc func(bindings.DataProxy[*repository.TagDBData], func())) func() {
	return func() {
		c.EditButton.Disable()

		editFunc(c.tag, func() {
			c.EditButton.Enable()
		})
	}
}

func (c *TagCard) deleteTag(callback func(*TagCard)) func() {
	logger := logging.With(c.deps.Logger(), keys.Component, "TagCard.deleteTag")
	return func() {
		ctx := context.Background()

		c.DeleteButton.Disable()

		tag, err := c.tag.Get()
		if err != nil {
			apperrors.Show(logger, c.parent, apperrors.Error(
				"Could not find tag data",
				apperrors.WithCause(err),
			), nil)
			return
		}

		logger := logging.With(logger, keys.TagID, tag.Tag.ID, keys.TagName, tag.Tag.Name) //nolint:govet // intentional

		conf := dialog.NewConfirm("Delete Tag?", fmt.Sprintf("Are you sure you want to delete the tag '%s'?", tag.Tag.Name), func(ok bool) {
			defer c.DeleteButton.Enable()
			if !ok {
				return
			}

			if err := c.deps.AppRepo().DeleteTag(ctx, c.TagID(), nil); err != nil {
				apperrors.Show(logger, c.parent, apperrors.Error(
					"Could not delete tag",
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

func (c *TagCard) Destroy() {}

func (c *TagCard) Layout(sz fyne.Size) {
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

func (c *TagCard) MinSize() fyne.Size {
	return fyne.Size{
		Height: 50,
		Width:  100,
	}
}

func (c *TagCard) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		c.ColorSwatch,
		c.NameLabel,
		c.EditButton,
		c.DeleteButton,
	}
}

func (c *TagCard) Refresh() {
	for _, o := range c.Objects() {
		o.Refresh()
	}
}
