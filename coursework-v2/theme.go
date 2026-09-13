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
	case theme.ColorNameMenuBackground, theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 31, G: 36, B: 41, A: 255}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 36, G: 41, B: 47, A: 255}
	case theme.ColorNameInputBorder, theme.ColorNameSeparator:
		return color.NRGBA{R: 78, G: 88, B: 96, A: 255}
	case theme.ColorNameButton, theme.ColorNamePrimary:
		return color.NRGBA{R: 122, G: 201, B: 106, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 235, G: 239, B: 235, A: 255}
	case theme.ColorNameForegroundOnPrimary:
		return color.NRGBA{R: 16, G: 24, B: 18, A: 255}
	case theme.ColorNamePlaceHolder, theme.ColorNameDisabled:
		return color.NRGBA{R: 164, G: 174, B: 180, A: 255}
	case theme.ColorNameHover:
		return color.NRGBA{R: 49, G: 60, B: 55, A: 255}
	case theme.ColorNameSelection, theme.ColorNameFocus:
		return color.NRGBA{R: 75, G: 130, B: 84, A: 255}
	case theme.ColorNameError:
		return color.NRGBA{R: 220, G: 90, B: 82, A: 255}
	case theme.ColorNameForegroundOnError:
		return color.NRGBA{R: 255, G: 245, B: 245, A: 255}
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 122, G: 201, B: 106, A: 255}
	case theme.ColorNameForegroundOnSuccess:
		return color.NRGBA{R: 16, G: 24, B: 18, A: 255}
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
