package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type courseworkTheme struct {
	base fyne.Theme
}

func (t courseworkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 22, G: 25, B: 29, A: 255}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 36, G: 41, B: 47, A: 255}
	case theme.ColorNameButton, theme.ColorNamePrimary:
		return color.NRGBA{R: 122, G: 201, B: 106, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 235, G: 239, B: 235, A: 255}
	case theme.ColorNameHeaderBackground:
		return color.NRGBA{R: 29, G: 34, B: 39, A: 255}
	}
	return t.base.Color(name, variant)
}

func (t courseworkTheme) Font(style fyne.TextStyle) fyne.Resource    { return t.base.Font(style) }
func (t courseworkTheme) Icon(name fyne.ThemeIconName) fyne.Resource { return t.base.Icon(name) }
func (t courseworkTheme) Size(name fyne.ThemeSizeName) float32       { return t.base.Size(name) }

func applyCourseworkTheme(myApp fyne.App) {
	myApp.Settings().SetTheme(courseworkTheme{base: theme.DefaultTheme()})
}
