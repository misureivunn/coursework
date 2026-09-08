package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type fstecTheme struct {
	base fyne.Theme
}

func (t fstecTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 245, G: 245, B: 245, A: 255}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 198, G: 45, B: 45, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 218, G: 222, B: 226, A: 255}
	case theme.ColorNameForegroundOnPrimary:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	case theme.ColorNameError:
		return color.NRGBA{R: 198, G: 45, B: 45, A: 255}
	case theme.ColorNameForegroundOnError:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	case theme.ColorNameHeaderBackground:
		return color.NRGBA{R: 232, G: 237, B: 242, A: 255}
	default:
		return t.base.Color(name, variant)
	}
}

func (t fstecTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.base.Font(style)
}

func (t fstecTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.base.Icon(name)
}

func (t fstecTheme) Size(name fyne.ThemeSizeName) float32 {
	return t.base.Size(name)
}

// ApplyFSTECTheme applies the common light FSTEC-inspired application palette.
func ApplyFSTECTheme(app fyne.App) {
	app.Settings().SetTheme(fstecTheme{base: theme.DefaultTheme()})
}
