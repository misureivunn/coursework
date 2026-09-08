package ui

import (
	"image/color"

	"szi-registry/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

func ShowLoginWindow(myApp fyne.App, db *gorm.DB) {
	myWindow := myApp.NewWindow("Вход в Реестр СЗИ")
	myWindow.Resize(fyne.NewSize(460, 360))

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Введите имя пользователя")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Введите пароль")

	errorLabel := widget.NewLabel("")
	errorLabel.Importance = widget.WarningImportance

	var loginButton *widget.Button

	loginFunc := func() {
		if usernameEntry.Text == "" || passwordEntry.Text == "" {
			errorLabel.SetText("Введите логин и пароль")
			return
		}
		user, err := services.AuthenticateUser(db, usernameEntry.Text, passwordEntry.Text)
		if err != nil {
			errorLabel.SetText("Неверный логин или пароль")
		} else {
			myWindow.Close()
			ShowSZIRegistryWindow(myApp, user.Username, db)
		}
	}

	loginButton = widget.NewButtonWithIcon("Войти", nil, loginFunc)
	loginButton.Importance = widget.DangerImportance

	buttonWithErrorContainer := container.NewVBox(loginButton, errorLabel)

	//регистраци
	registerButton := widget.NewButton("Регистрация", func() {
		ShowRegisterWindow(myApp, db)
	})
	registerButton.Importance = widget.DangerImportance

	usernameEntry.Resize(fyne.NewSize(400, 36))
	passwordEntry.Resize(fyne.NewSize(400, 36))

	title := canvas.NewText("Вход в реестр СЗИ", color.NRGBA{R: 26, G: 60, B: 110, A: 255})
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	content := container.NewVBox(
		title,
		widget.NewSeparator(), // Разделитель
		container.NewPadded(container.NewVBox(
			widget.NewLabel("Имя пользователя"),
			usernameEntry,
			widget.NewLabel("Пароль"),
			passwordEntry,
		)), // Форма с отступами
		container.NewPadded(container.NewVBox(
			buttonWithErrorContainer,
			registerButton,
		)),
	)

	background := canvas.NewRectangle(color.RGBA{R: 240, G: 240, B: 240, A: 255})
	backgroundContainer := container.NewStack(background, container.NewCenter(container.NewPadded(content)))

	myWindow.SetContent(backgroundContainer)
	myWindow.Show()
}
