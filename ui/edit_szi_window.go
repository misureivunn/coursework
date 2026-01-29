package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
	"szi-registry/models"
	"szi-registry/services"
	"szi-registry/utils/ui"
)

func ShowEditSziWindow(myApp fyne.App, record models.SZIRecord, db *gorm.DB) {
	myWindow := myApp.NewWindow("Редактировать СЗИ")
	myWindow.Resize(fyne.NewSize(600, 500))

	nameEntry := widget.NewEntry()
	nameEntry.SetText(record.Name)
	// Заменяем текстовое поле на выпадающий список для типа СЗИ
	typeSelector := widget.NewSelect([]string{"СКЗИ", "ОС", "СЗИ КС", "СЗИ СКЗИ", "Другое"}, func(value string) {})
	typeSelector.SetSelected(record.Type)
	certNumberEntry := widget.NewEntry()
	certNumberEntry.SetText(record.CertNumber)

	issueDateEntry := widget.NewEntry()
	issueDateEntry.SetText(record.CertIssueDate.Format("2006-01-02"))
	issueDateEntry.PlaceHolder = "ГГГГ-ММ-ДД (например: 2025-01-02)"

	expiryDateEntry := widget.NewEntry()
	expiryDateEntry.SetText(record.CertExpiryDate.Format("2006-01-02"))
	expiryDateEntry.PlaceHolder = "ГГГГ-ММ-ДД (например: 2026-01-02)"

	locationEntry := widget.NewEntry()
	locationEntry.SetText(record.Location)
	statusEntry := widget.NewEntry()
	statusEntry.SetText(record.Status)
	manufacturerEntry := widget.NewEntry()
	manufacturerEntry.SetText(record.Manufacturer)
	softwareVersionEntry := widget.NewEntry()
	softwareVersionEntry.SetText(record.SoftwareVersion)

	// Добавляем новые поля для классификации СЗИ от НСД
	purposeSelector := widget.NewSelect([]string{"АС", "ИВК", "Универсальное"}, func(value string) {})
	purposeSelector.SetSelected(record.Purpose)
	deploymentTypeSelector := widget.NewSelect([]string{"Клиент-сервер", "Автономное", "АПК", "Виртуальное"}, func(value string) {})
	deploymentTypeSelector.SetSelected(record.DeploymentType)
	classProtectionSelector := widget.NewSelect([]string{"1", "2", "3А", "4", "5"}, func(value string) {})
	classProtectionSelector.SetSelected(record.ClassProtection)

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
			{Text: "Статус", Widget: statusEntry},
			{Text: "Производитель", Widget: manufacturerEntry},
			{Text: "Версия ПО", Widget: softwareVersionEntry},
		},
		OnSubmit: func() {
			issueDate, err := ui.ParseDate(issueDateEntry.Text)
			if err != nil {
				ui.ShowErrorDialog("Неверный формат даты выдачи. Используйте ГГГГ-ММ-ДД", myWindow.Canvas())
				return
			}

			expiryDate, err := ui.ParseDate(expiryDateEntry.Text)
			if err != nil {
				ui.ShowErrorDialog("Неверный формат срока действия. Используйте ГГГГ-ММ-ДД", myWindow.Canvas())
				return
			}

			record.Name = nameEntry.Text
			record.Type = typeSelector.Selected
			record.CertNumber = certNumberEntry.Text
			record.CertIssueDate = issueDate
			record.CertExpiryDate = expiryDate
			record.Location = locationEntry.Text
			record.Status = statusEntry.Text
			record.Manufacturer = manufacturerEntry.Text
			record.SoftwareVersion = softwareVersionEntry.Text

			// Обновляем новые поля для классификации СЗИ от НСД
			record.Purpose = purposeSelector.Selected
			record.DeploymentType = deploymentTypeSelector.Selected
			record.ClassProtection = classProtectionSelector.Selected

			err = services.UpdateSziRecord(db, &record)
			if err != nil {
				dialog := widget.NewModalPopUp(widget.NewLabel("Ошибка при обновлении записи: "+err.Error()), myWindow.Canvas())
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
	// typeSelector - это выпадающий список, ему не нужен метод Resize
	certNumberEntry.Resize(fyne.NewSize(350, 30))
	issueDateEntry.Resize(fyne.NewSize(350, 30))
	expiryDateEntry.Resize(fyne.NewSize(350, 30))
	locationEntry.Resize(fyne.NewSize(350, 30))
	statusEntry.Resize(fyne.NewSize(350, 30))
	manufacturerEntry.Resize(fyne.NewSize(350, 30))
	softwareVersionEntry.Resize(fyne.NewSize(350, 30))

	// Улучшенное оформление формы
	scrollContainer := container.NewVScroll(form)
	scrollContainer.SetMinSize(fyne.NewSize(580, 400))
	buttonContainer := container.NewHBox(
		widget.NewButton("Назад", func() { myWindow.Close() }),
		widget.NewButton("Сохранить", func() { form.OnSubmit() }),
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
