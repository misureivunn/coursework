package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func showAppError(message string, window fyne.Window) {
	dialog.ShowError(fmt.Errorf("%s", message), window)
}
