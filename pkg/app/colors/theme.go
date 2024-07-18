package colors

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type AppTheme struct{}

var _ fyne.Theme = (*AppTheme)(nil)

var ColorNameInvertedForeground fyne.ThemeColorName = "invertedForeground"

func (t *AppTheme) Color(cn fyne.ThemeColorName, tv fyne.ThemeVariant) color.Color {
	if cn == ColorNameInvertedForeground {
		cn = theme.ColorNameForeground
		switch tv {
		case theme.VariantDark:
			return theme.DefaultTheme().Color(cn, theme.VariantLight)
		case theme.VariantLight:
			return theme.DefaultTheme().Color(cn, theme.VariantDark)
		}

	}
	return theme.DefaultTheme().Color(cn, tv)
}

func (t *AppTheme) Font(ts fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(ts)
}

func (t *AppTheme) Icon(i fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(i)
}

func (t *AppTheme) Size(sn fyne.ThemeSizeName) float32 {
	if sn == theme.SizeNameCaptionText {
		return theme.DefaultTheme().Size(theme.SizeNameText) * 2 / 3
	}

	return theme.DefaultTheme().Size(sn)
}
