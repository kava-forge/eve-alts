package app

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"
)

type MiniTag struct {
	widget.BaseWidget

	logger       logging.Logger
	NameLabel    *widget.RichText
	ColorSwatch  *ColorSwatch
	Dimmed       bool
	UnDimmedBold bool

	size fyne.ThemeSizeName
}

func NewMiniTag(logger logging.Logger, text string, c color.Color, size fyne.ThemeSizeName) *MiniTag {
	cc := &MiniTag{
		logger:      logger,
		NameLabel:   widget.NewRichTextWithText(text),
		ColorSwatch: NewColorSwatch(logger, c),
		size:        size,
	}
	cc.refreshStyle()
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

func normalizeC(c, a uint32) float64 {
	rc := float64(c) / float64(a)
	if rc <= 0.04045 {
		rc = rc / 12.92
	} else {
		rc = math.Exp(math.Log((rc+0.055)/1.055) * 2.4)
		// rc = math.Pow(, 2.4)
	}

	return rc
}

// https://stackoverflow.com/a/3943023
func UseDarkText(bgc color.Color) bool {
	r, g, b, a := bgc.RGBA()

	rc, gc, bc := normalizeC(r, a), normalizeC(g, a), normalizeC(b, a)
	l := 0.2126*rc + 0.7152*gc + 0.0722*bc

	return l > 0.179
}

func (c *MiniTag) refreshStyle() {
	defer c.NameLabel.Refresh()

	c.NameLabel.Wrapping = fyne.TextWrapOff

	textColorName := theme.ColorNameForeground

	if !c.Dimmed {
		var darkText, lightText fyne.ThemeColorName
		if fyne.CurrentApp().Settings().ThemeVariant() == theme.VariantLight {
			darkText = theme.ColorNameForeground
			lightText = ColorNameInvertedForeground
		} else {
			lightText = theme.ColorNameForeground
			darkText = ColorNameInvertedForeground
		}

		if UseDarkText(c.ColorSwatch.Color()) {
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
	c.refreshStyle()
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
