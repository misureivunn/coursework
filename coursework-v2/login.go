package main

import (
	"fmt"
	"image/color"

	"szi-registry/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

func showLogin(myApp fyne.App, db *gorm.DB) {
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
		user, err := services.AuthenticateUser(db, username.Text, password.Text)
		if err != nil {
			message.SetText("Проверьте имя пользователя и пароль")
			return
		}
		window.Close()
		showWorkspace(myApp, db, user)
	})
	login.Importance = widget.HighImportance

	register := widget.NewButton("Создать аккаунт", func() {
		showRegister(myApp, db, window)
	})

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

func showRegister(myApp fyne.App, db *gorm.DB, loginWindow fyne.Window) {
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
		if password.Text != confirm.Text {
			dialog.ShowError(fmt.Errorf("пароли не совпадают"), window)
			return
		}
		if _, err := services.CreateUser(db, username.Text, password.Text); err != nil {
			dialog.ShowError(err, window)
			return
		}
		window.Close()
		dialog.ShowInformation("Готово", "Аккаунт создан", loginWindow)
	}
	window.SetContent(container.NewBorder(nil, widget.NewButton("Закрыть", func() { window.Close() }), nil, nil, form))
	window.Show()
}
