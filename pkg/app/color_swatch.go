package app

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"
)

type ColorSwatch struct {
	widget.BaseWidget

	logger logging.Logger
	rect   *canvas.Rectangle
}

func NewColorSwatch(logger logging.Logger, c color.Color) *ColorSwatch {
	s := &ColorSwatch{
		logger: logger,
		rect:   canvas.NewRectangle(c),
	}
	s.ExtendBaseWidget(s)

	return s
}

func (s *ColorSwatch) SetCornerRadius(cr float32) {
	s.rect.CornerRadius = cr
}

func (s *ColorSwatch) Color() color.Color {
	return s.rect.FillColor
}

func (s *ColorSwatch) SetColor(c color.Color) {
	defer s.Refresh()

	s.rect.FillColor = c
}

func (s *ColorSwatch) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(s.rect)
}

type TappableColorSwatch struct {
	ColorSwatch
	OnTapped func(*fyne.PointEvent)
}

func NewTappableColorSwatch(logger logging.Logger, c color.Color) *TappableColorSwatch {
	return &TappableColorSwatch{
		ColorSwatch: *NewColorSwatch(logger, c),
	}
}

func (s *TappableColorSwatch) Tapped(e *fyne.PointEvent) {
	level.Debug(s.logger).Message("tapped ColorSwatch")
	if s.OnTapped != nil {
		s.OnTapped(e)
	}
}
