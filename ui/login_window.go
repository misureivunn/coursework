package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
	"szi-registry/services"
)

func ShowLoginWindow(myApp fyne.App, db *gorm.DB) {
	myWindow := myApp.NewWindow("Вход в Реестр СЗИ")
	myWindow.Resize(fyne.NewSize(400, 250))

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Введите имя пользователя")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Введите пароль")

	// Label для отображения ошибок
	errorLabel := widget.NewLabel("")
	errorLabel.Importance = widget.WarningImportance

	// Define the login button later after all variables are declared
	var loginButton *widget.Button

	loginFunc := func() {
		// Real authentication - check credentials against DB
		user, err := services.AuthenticateUser(db, usernameEntry.Text, passwordEntry.Text)
		if err != nil {
			// Authentication failed
			errorLabel.SetText(err.Error())
		} else {
			// Authentication successful - show main window with user info
			myWindow.Close()
			ShowMainWindow(myApp, user.Username, db)
		}
	}

	loginButton = widget.NewButtonWithIcon("Войти", nil, loginFunc)

	// Контейнер для кнопки и сообщения об ошибке
	buttonWithErrorContainer := container.NewVBox(loginButton, errorLabel)

	// Кнопка регистрации
	registerButton := widget.NewButton("Регистрация", func() {
		ShowRegisterWindow(myApp, db)
	})

	// Устанавливаем минимальный размер для полей ввода
	usernameEntry.Resize(fyne.NewSize(400, 30))
	passwordEntry.Resize(fyne.NewSize(400, 30))

	// Центрируем содержимое
	content := container.NewVBox(
		widget.NewLabel("Добро пожаловать в Реестр СЗИ"), // Заголовок
		widget.NewSeparator(), // Разделитель
		container.NewPadded(container.NewVBox(
			widget.NewLabel("Имя пользователя"),
			usernameEntry,
			widget.NewLabel("Пароль"),
			passwordEntry,
		)), // Форма с отступами
		container.NewPadded(buttonWithErrorContainer), // Кнопка с отступами
		container.NewPadded(registerButton),           // Кнопка регистрации
	)

	// Добавим фоновый цвет для улучшения визуального восприятия
	background := canvas.NewRectangle(color.RGBA{R: 240, G: 240, B: 240, A: 255})
	backgroundContainer := container.NewStack(background, container.NewCenter(content))

	myWindow.SetContent(backgroundContainer)
	myWindow.Show()
}
