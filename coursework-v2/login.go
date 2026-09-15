package main

import (
	"fmt"
	"image/color"

	"coursework-v2/internal/application"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func showLogin(myApp fyne.App, app application.App) {
	window := myApp.NewWindow("Реестр СЗИ")
	window.Resize(fyne.NewSize(460, 390))

	username := widget.NewEntry()
	username.SetPlaceHolder("Имя пользователя")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Пароль")
	message := widget.NewLabel("")

	title := canvas.NewText("РЕЕСТР СЗИ", color.NRGBA{R: 122, G: 201, B: 106, A: 255})
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}

	login := widget.NewButton("Войти", func() {
		if username.Text == "" || password.Text == "" {
			message.SetText("Введите имя пользователя и пароль")
			return
		}
		user, err := app.Login(username.Text, password.Text)
		if err != nil {
			message.SetText("Вход не выполнен. Проверьте имя пользователя и пароль")
			return
		}
		window.Close()
		showWorkspace(myApp, app, user)
	})
	login.Importance = widget.HighImportance

	register := widget.NewButton("Создать аккаунт", func() {
		showRegister(myApp, app, window)
	})
	register.Importance = widget.HighImportance

	form := container.NewVBox(
		container.NewCenter(title),
		widget.NewSeparator(),
		widget.NewLabel("Вход в рабочее пространство"),
		username,
		password,
		login,
		register,
		message,
	)
	window.SetContent(container.NewCenter(container.NewPadded(form)))
	window.Show()
}

func showRegister(myApp fyne.App, app application.App, loginWindow fyne.Window) {
	window := myApp.NewWindow("Новый пользователь")
	window.Resize(fyne.NewSize(420, 300))
	username := widget.NewEntry()
	username.SetPlaceHolder("Имя пользователя")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Пароль")
	confirm := widget.NewPasswordEntry()
	confirm.SetPlaceHolder("Повторите пароль")

	form := widget.NewForm(
		widget.NewFormItem("Логин", username),
		widget.NewFormItem("Пароль", password),
		widget.NewFormItem("Повтор", confirm),
	)
	form.OnSubmit = func() {
		if username.Text == "" || password.Text == "" || confirm.Text == "" {
			dialog.ShowError(fmt.Errorf("Заполните логин, пароль и подтверждение пароля"), window)
			return
		}
		if password.Text != confirm.Text {
			dialog.ShowError(fmt.Errorf("Пароли не совпадают. Повторите ввод"), window)
			return
		}
		if err := app.Register(username.Text, password.Text); err != nil {
			dialog.ShowError(fmt.Errorf("Не удалось создать аккаунт: %v", err), window)
			return
		}
		window.Close()
		dialog.ShowInformation("Аккаунт создан", "Регистрация завершена. Введите данные в окне входа", loginWindow)
	}
	closeButton := widget.NewButton("Закрыть", func() { window.Close() })
	closeButton.Importance = widget.HighImportance
	window.SetContent(container.NewBorder(nil, closeButton, nil, nil, form))
	window.Show()
}
