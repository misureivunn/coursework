package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
	"szi-registry/models"
	"szi-registry/services"
)

func ShowRegisterWindow(myApp fyne.App, db *gorm.DB) {
	MyWindow := myApp.NewWindow("Регистрация")
	MyWindow.Resize(fyne.NewSize(500, 400))

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Введите имя пользователя")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Введите пароль")

	confirmPasswordEntry := widget.NewPasswordEntry()
	confirmPasswordEntry.SetPlaceHolder("Подтвердите пароль")

	errorLabel := widget.NewLabel("")
	errorLabel.Importance = widget.WarningImportance

	// Создаем контейнер для формы с вертикальной компоновкой, где каждая метка над своим полем
	formContainer := container.NewVBox(
		container.NewPadded(container.NewVBox(
			widget.NewLabel("Имя пользователя"),
			usernameEntry,
		)),
		container.NewPadded(container.NewVBox(
			widget.NewLabel("Пароль"),
			passwordEntry,
		)),
		container.NewPadded(container.NewVBox(
			widget.NewLabel("Подтверждение пароля"),
			confirmPasswordEntry,
		)),
	)

	// Создаем функцию отправки формы
	onSubmit := func() {
		// Проверяем, совпадают ли пароли
		if passwordEntry.Text != confirmPasswordEntry.Text {
			errorLabel.SetText("Пароли не совпадают")
			return
		}

		// Создаем пользователя
		user := &models.User{
			Username: usernameEntry.Text,
			Role:     "user",
		}

		// Создаем пользователя через сервис
		_, err := services.CreateUser(db, user.Username, passwordEntry.Text)
		if err != nil {
			errorLabel.SetText("Ошибка при создании пользователя: " + err.Error())
			return
		}

		// Закрываем окно регистрации
		MyWindow.Close()
	}

	// Устанавливаем минимальный размер для полей ввода
	usernameEntry.Resize(fyne.NewSize(400, 30))
	passwordEntry.Resize(fyne.NewSize(400, 30))
	confirmPasswordEntry.Resize(fyne.NewSize(400, 30))

	// Кнопки
	buttonContainer := container.NewHBox(
		widget.NewButton("Назад", func() { MyWindow.Close() }),
		widget.NewButton("Зарегистрироваться", func() { onSubmit() }),
	)

	content := container.NewBorder(
		nil,
		buttonContainer,
		nil,
		nil,
		formContainer,
	)

	mainContent := container.NewVBox(
		widget.NewLabel("Регистрация нового пользователя"),
		widget.NewSeparator(),
		content,
		errorLabel,
	)

	// Добавим фоновый цвет для улучшения визуального восприятия
	background := canvas.NewRectangle(color.RGBA{R: 240, G: 240, B: 240, A: 255})
	backgroundContainer := container.NewStack(background, container.NewCenter(mainContent))

	MyWindow.SetContent(backgroundContainer)
	MyWindow.Show()
}
