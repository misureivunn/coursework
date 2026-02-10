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
	myWindow.Resize(fyne.NewSize(400, 250))

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Введите имя пользователя")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Введите пароль")

	errorLabel := widget.NewLabel("")
	errorLabel.Importance = widget.WarningImportance

	var loginButton *widget.Button

	loginFunc := func() {
		user, err := services.AuthenticateUser(db, usernameEntry.Text, passwordEntry.Text)
		if err != nil {
			errorLabel.SetText(err.Error())
		} else {
			myWindow.Close()
			ShowDashboardWindow(myApp, user.Username, db)
		}
	}

	loginButton = widget.NewButtonWithIcon("Войти", nil, loginFunc)

	buttonWithErrorContainer := container.NewVBox(loginButton, errorLabel)

	//регистраци
	registerButton := widget.NewButton("Регистрация", func() {
		ShowRegisterWindow(myApp, db)
	})

	usernameEntry.Resize(fyne.NewSize(400, 30))
	passwordEntry.Resize(fyne.NewSize(400, 30))

	content := container.NewVBox(
		widget.NewLabel("Добро пожаловать в Реестр СЗИ"),
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
	backgroundContainer := container.NewStack(background, container.NewCenter(content))

	myWindow.SetContent(backgroundContainer)
	myWindow.Show()
}
