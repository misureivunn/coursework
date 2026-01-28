package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
	"szi-registry/models"
	"szi-registry/services"
)

func ShowAddSziWindow(myApp fyne.App, userID uint, db *gorm.DB) {
	myWindow := myApp.NewWindow("Добавить СЗИ")
	myWindow.Resize(fyne.NewSize(600, 500))

	nameEntry := widget.NewEntry()
	typeEntry := widget.NewEntry()
	certNumberEntry := widget.NewEntry()

	issueDateEntry := widget.NewEntry()
	issueDateEntry.PlaceHolder = "ГГГГ-ММ-ДД (например: 2025-01-02)"

	expiryDateEntry := widget.NewEntry()
	expiryDateEntry.PlaceHolder = "ГГГГ-ММ-ДД (например: 2026-01-02)"

	locationEntry := widget.NewEntry()
	manufacturerEntry := widget.NewEntry()
	softwareVersionEntry := widget.NewEntry()
	contactPersonEntry := widget.NewEntry()
	documentationLinkEntry := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Наименование СЗИ", Widget: nameEntry},
			{Text: "Тип СЗИ", Widget: typeEntry},
			{Text: "Номер сертификата", Widget: certNumberEntry},
			{Text: "Дата выдачи (ГГГГ-ММ-ДД)", Widget: issueDateEntry},
			{Text: "Срок действия (ГГГГ-ММ-ДД)", Widget: expiryDateEntry},
			{Text: "Место установки", Widget: locationEntry},
			{Text: "Производитель", Widget: manufacturerEntry},
			{Text: "Версия ПО", Widget: softwareVersionEntry},
			{Text: "Контактное лицо", Widget: contactPersonEntry},
			{Text: "Ссылка на документацию", Widget: documentationLinkEntry},
		},
		OnSubmit: func() {
			issueDate, err := time.Parse("2006-01-02", issueDateEntry.Text)
			if err != nil {
				dialog := widget.NewModalPopUp(widget.NewLabel("Неверный формат даты выдачи. Используйте ГГГГ-ММ-ДД"), myWindow.Canvas())
				dialog.Show()
				return
			}

			expiryDate, err := time.Parse("2006-01-02", expiryDateEntry.Text)
			if err != nil {
				dialog := widget.NewModalPopUp(widget.NewLabel("Неверный формат срока действия. Используйте ГГГГ-ММ-ДД"), myWindow.Canvas())
				dialog.Show()
				return
			}

			newRecord := &models.SZIRecord{
				Name:              nameEntry.Text,
				Type:              typeEntry.Text,
				CertNumber:        certNumberEntry.Text,
				CertIssueDate:     issueDate,
				CertExpiryDate:    expiryDate,
				Location:          locationEntry.Text,
				UserID:            userID, // Используем переданный ID пользователя
				Manufacturer:      manufacturerEntry.Text,
				SoftwareVersion:   softwareVersionEntry.Text,
				ContactPerson:     contactPersonEntry.Text,
				DocumentationLink: documentationLinkEntry.Text,
			}

			if expiryDate.Before(time.Now()) {
				newRecord.Status = "Просрочено"
			} else {
				newRecord.Status = "Актуально"
			}

			err = services.AddSziRecord(db, newRecord)
			if err != nil {
				dialog := widget.NewModalPopUp(widget.NewLabel("Ошибка при добавлении записи: "+err.Error()), myWindow.Canvas())
				dialog.Show()
				return
			}

			myWindow.Close()
		},
		OnCancel: func() {
			myWindow.Close()
		},
	}

	// Устанавливаем минимальный размер для полей ввода
	nameEntry.Resize(fyne.NewSize(350, 30))
	typeEntry.Resize(fyne.NewSize(350, 30))
	certNumberEntry.Resize(fyne.NewSize(350, 30))
	issueDateEntry.Resize(fyne.NewSize(350, 30))
	expiryDateEntry.Resize(fyne.NewSize(350, 30))
	locationEntry.Resize(fyne.NewSize(350, 30))
	manufacturerEntry.Resize(fyne.NewSize(350, 30))
	softwareVersionEntry.Resize(fyne.NewSize(350, 30))
	contactPersonEntry.Resize(fyne.NewSize(350, 30))
	documentationLinkEntry.Resize(fyne.NewSize(350, 30))

	// Улучшенное оформление формы
	scrollContainer := container.NewVScroll(form)
	scrollContainer.SetMinSize(fyne.NewSize(580, 400))
	buttonContainer := container.NewHBox(
		widget.NewButton("Назад", func() { myWindow.Close() }),
		widget.NewButton("Добавить", func() { form.OnSubmit() }),
	)

	content := container.NewBorder(
		nil,
		buttonContainer,
		nil,
		nil,
		scrollContainer,
	)

	myWindow.SetContent(content)
	myWindow.Show()
}
