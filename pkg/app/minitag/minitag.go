package minitag

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"
	"github.com/kava-forge/eve-alts/pkg/app/colors"
)

type MiniTag struct {
	widget.BaseWidget

	logger       logging.Logger
	NameLabel    *widget.RichText
	ColorSwatch  *colors.ColorSwatch
	Dimmed       bool
	UnDimmedBold bool

	size fyne.ThemeSizeName
}

func New(logger logging.Logger, text string, c color.Color, size fyne.ThemeSizeName) *MiniTag {
	cc := &MiniTag{
		logger:      logger,
		NameLabel:   widget.NewRichTextWithText(text),
		ColorSwatch: colors.NewColorSwatch(logger, c),
		size:        size,
	}
	cc.RefreshStyle()
	cc.ColorSwatch.SetCornerRadius(theme.InnerPadding() / 2)
	cc.ExtendBaseWidget(cc)

	return cc
}

func (c *MiniTag) CreateRenderer() fyne.WidgetRenderer {
	return c
}

func (c *MiniTag) SetText(text string) {
	c.SetTextAt(0, text)
}

func (c *MiniTag) SetTextAt(idx int, text string) {
	c.NameLabel.Segments[idx].(*widget.TextSegment).Text = text
	c.NameLabel.Refresh()
}

func (c *MiniTag) RefreshStyle() {
	defer c.NameLabel.Refresh()

	c.NameLabel.Wrapping = fyne.TextWrapOff

	textColorName := theme.ColorNameForeground

	if !c.Dimmed {
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
	}

	for _, segi := range c.NameLabel.Segments {
		if seg, ok := segi.(*widget.TextSegment); ok {
			seg.Style = widget.RichTextStyle{
				Alignment: fyne.TextAlignLeading,
				ColorName: textColorName,
				Inline:    true,
				SizeName:  c.size,
				TextStyle: fyne.TextStyle{},
			}
			if c.Dimmed {
				seg.Style.TextStyle.Bold = false
			} else {
				seg.Style.TextStyle.Bold = c.UnDimmedBold
			}
		}
		segi.Visual().Refresh()
	}
}

// The WidgetRenderer interface

func (c *MiniTag) Destroy() {}

func (c *MiniTag) Layout(sz fyne.Size) {
	level.Debug(c.logger).Message("MiniTag.Layout", "size", sz, "text", c.NameLabel.String())

	c.NameLabel.Refresh()
	c.RefreshStyle()
	c.NameLabel.Move(fyne.Position{X: -theme.InnerPadding() / 2, Y: -2 * theme.InnerPadding() / 3})

	// half-size paddings
	c.ColorSwatch.Resize(sz)
	c.ColorSwatch.Move(fyne.Position{X: 0, Y: 0})

	if c.Dimmed {
		c.ColorSwatch.Hide()
	} else {
		c.ColorSwatch.Show()
	}
}

func (c *MiniTag) MinSize() fyne.Size {
	// fontSize := fyne.CurrentApp().Settings().Theme().Size(c.size)
	// tsz := fyne.MeasureText(c.NameLabel.String(), fontSize, c.NameLabel.Segments[0].(*widget.TextSegment).Style.TextStyle)
	tsz := c.NameLabel.MinSize()

	// half-size padding H, 1/3 padding V
	return fyne.Size{Width: tsz.Width - theme.InnerPadding(), Height: tsz.Height - 4*theme.InnerPadding()/3}
}

func (c *MiniTag) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		c.ColorSwatch,
		c.NameLabel,
	}
}

func (c *MiniTag) Refresh() {
	c.Layout(c.MinSize())
	for _, o := range c.Objects() {
		o.Refresh()
	}
}
