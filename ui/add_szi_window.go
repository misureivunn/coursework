package ui

import (
	"szi-registry/models"
	"szi-registry/services"
	"szi-registry/utils/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

func ShowAddSziWindow(myApp fyne.App, userID uint, db *gorm.DB) {
	myWindow := myApp.NewWindow("Добавить СЗИ")
	myWindow.Resize(fyne.NewSize(600, 500))

	nameEntry := widget.NewEntry()
	// Заменяем текстовое поле на выпадающий список для типа СЗИ
	typeSelector := widget.NewSelect([]string{"СКЗИ", "ОС", "СЗИ КС", "СЗИ СКЗИ", "Другое"}, func(value string) {})
	certNumberEntry := widget.NewEntry()

	issueDateEntry := widget.NewEntry()
	issueDateEntry.PlaceHolder = "ГГГГ-ММ-ДД (например: 2025-01-02)"

	expiryDateEntry := widget.NewEntry()
	expiryDateEntry.PlaceHolder = "ГГГГ-ММ-ДД (например: 2026-01-02)"

	locationEntry := widget.NewEntry()
	manufacturerEntry := widget.NewEntry()
	softwareEntry := widget.NewEntry()
	notesEntry := widget.NewMultiLineEntry()

	// Добавляем новые поля для классификации СЗИ от НСД
	purposeSelector := widget.NewSelect([]string{"АС", "ИВК", "Универсальное"}, func(value string) {})
	deploymentTypeSelector := widget.NewSelect([]string{"Клиент-сервер", "Автономное", "АПК", "Виртуальное"}, func(value string) {})
	classProtectionSelector := widget.NewSelect([]string{"1", "2", "3А", "4", "5"}, func(value string) {})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Наименование СЗИ", Widget: nameEntry},
			{Text: "Тип СЗИ", Widget: typeSelector},
			{Text: "Назначение", Widget: purposeSelector},
			{Text: "Тип развертывания", Widget: deploymentTypeSelector},
			{Text: "Класс защищенности", Widget: classProtectionSelector},
			{Text: "Номер сертификата", Widget: certNumberEntry},
			{Text: "Дата выдачи (ГГГГ-ММ-ДД)", Widget: issueDateEntry},
			{Text: "Срок действия (ГГГГ-ММ-ДД)", Widget: expiryDateEntry},
			{Text: "Место установки", Widget: locationEntry},
			{Text: "Производитель", Widget: manufacturerEntry},
			{Text: "Версия ПО", Widget: softwareEntry},
			{Text: "Примечания", Widget: notesEntry},
		},
		OnSubmit: func() {
			issueDate, err := ui.ParseDate(issueDateEntry.Text)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}

			expiryDate, err := ui.ParseDate(expiryDateEntry.Text)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}

			newRecord := &models.SZIRecord{
				Name:            nameEntry.Text,
				Type:            typeSelector.Selected, // Используем выбранный тип из выпадающего списка
				CertNumber:      certNumberEntry.Text,
				CertIssueDate:   issueDate,
				CertExpiryDate:  expiryDate,
				Location:        locationEntry.Text,
				UserID:          userID,
				Manufacturer:    manufacturerEntry.Text,
				SoftwareVersion: softwareEntry.Text,
				Notes:           notesEntry.Text,
				Purpose:         purposeSelector.Selected,
				DeploymentType:  deploymentTypeSelector.Selected,
				ClassProtection: classProtectionSelector.Selected,
			}

			err = services.AddSziRecord(db, newRecord)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}

			myWindow.Close()
		},
		OnCancel: func() {
			myWindow.Close()
		},
	}

	nameEntry.Resize(fyne.NewSize(350, 30))
	certNumberEntry.Resize(fyne.NewSize(350, 30))
	issueDateEntry.Resize(fyne.NewSize(350, 30))
	expiryDateEntry.Resize(fyne.NewSize(350, 30))
	locationEntry.Resize(fyne.NewSize(350, 30))
	manufacturerEntry.Resize(fyne.NewSize(350, 30))
	softwareEntry.Resize(fyne.NewSize(350, 30))

	scrollContainer := container.NewVScroll(form)
	scrollContainer.SetMinSize(fyne.NewSize(580, 400))
	buttonContainer := container.NewHBox(
		widget.NewButton("Назад", func() { myWindow.Close() }),
	)
	addButton := widget.NewButton("Добавить", func() { form.OnSubmit() })
	addButton.Importance = widget.DangerImportance
	buttonContainer.Add(addButton)

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
