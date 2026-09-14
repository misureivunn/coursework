package main

import (
	"fmt"

	"szi-registry/models"
	"szi-registry/services"
	appui "szi-registry/utils/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

func showRecordForm(myApp fyne.App, db *gorm.DB, userID uint, record *models.SZIRecord, onClosed func()) {
	isEdit := record != nil
	title := "Добавить средство защиты"
	if isEdit {
		title = "Редактировать средство защиты"
	}
	window := myApp.NewWindow(title)
	window.Resize(fyne.NewSize(720, 700))

	name := widget.NewEntry()
	typeSelect := widget.NewSelect(models.SZITypeValues, nil)
	protectionClass := widget.NewSelect(models.ProtectionClassValues, nil)
	trustLevel := widget.NewSelect(models.TrustLevelValues, nil)
	acClass := widget.NewSelect(models.ACClassValues, nil)
	vendor := widget.NewEntry()
	version := widget.NewEntry()
	certificateNumber := widget.NewEntry()
	scheme := widget.NewSelect(models.CertificationSchemeValues, nil)
	issueDate := widget.NewEntry()
	expiryDate := widget.NewEntry()
	installLocation := widget.NewEntry()
	responsiblePerson := widget.NewEntry()
	notes := widget.NewMultiLineEntry()

	issueDate.SetPlaceHolder("ГГГГ-ММ-ДД")
	expiryDate.SetPlaceHolder("ГГГГ-ММ-ДД")
	if isEdit {
		name.SetText(record.Name)
		typeSelect.SetSelected(valueOrFallback(record.SZIType, record.Type))
		protectionClass.SetSelected(valueOrFallback(record.ProtectionClass, record.ClassProtection))
		trustLevel.SetSelected(record.TrustLevel)
		acClass.SetSelected(record.ACClass)
		vendor.SetText(valueOrFallback(record.Vendor, record.Manufacturer))
		version.SetText(valueOrFallback(record.Version, record.SoftwareVersion))
		certificateNumber.SetText(valueOrFallback(record.CertificateNumber, record.CertNumber))
		scheme.SetSelected(record.CertificationScheme)
		issueDate.SetText(formatReferenceDate(record.IssueDate, record.CertIssueDate))
		expiryDate.SetText(formatReferenceDate(record.ExpiryDate, record.CertExpiryDate))
		installLocation.SetText(valueOrFallback(record.InstallLocation, record.Location))
		responsiblePerson.SetText(record.ResponsiblePerson)
		notes.SetText(record.Notes)
	}

	form := widget.NewForm(
		widget.NewFormItem("Наименование СЗИ", name),
		widget.NewFormItem("Тип СЗИ", typeSelect),
		widget.NewFormItem("Класс защищённости СВТ", protectionClass),
		widget.NewFormItem("Уровень доверия", trustLevel),
		widget.NewFormItem("Класс АС", acClass),
		widget.NewFormItem("Заявитель / разработчик", vendor),
		widget.NewFormItem("Версия ПО / исполнение", version),
		widget.NewFormItem("Номер сертификата", certificateNumber),
		widget.NewFormItem("Схема сертификации", scheme),
		widget.NewFormItem("Дата выдачи", issueDate),
		widget.NewFormItem("Срок действия", expiryDate),
		widget.NewFormItem("Место установки", installLocation),
		widget.NewFormItem("Ответственный сотрудник", responsiblePerson),
		widget.NewFormItem("Примечания", notes),
	)

	form.OnSubmit = func() {
		if name.Text == "" || typeSelect.Selected == "" || certificateNumber.Text == "" {
			dialog.ShowError(fmt.Errorf("заполните наименование, тип СЗИ и номер сертификата"), window)
			return
		}
		parsedIssue, err := appui.ParseDate(issueDate.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("дата выдачи: %v", err), window)
			return
		}
		parsedExpiry, err := appui.ParseDate(expiryDate.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("срок действия: %v", err), window)
			return
		}
		if record == nil {
			record = &models.SZIRecord{UserID: userID}
		}
		record.OwnerID = record.UserID
		record.Name = name.Text
		record.Type = typeSelect.Selected
		record.CertNumber = certificateNumber.Text
		record.CertIssueDate = parsedIssue
		record.CertExpiryDate = parsedExpiry
		record.Location = installLocation.Text
		record.Manufacturer = vendor.Text
		record.SoftwareVersion = version.Text
		record.ClassProtection = protectionClass.Selected
		record.SZIType = typeSelect.Selected
		record.ProtectionClass = protectionClass.Selected
		record.TrustLevel = trustLevel.Selected
		record.ACClass = acClass.Selected
		record.Vendor = vendor.Text
		record.Version = version.Text
		record.CertificateNumber = certificateNumber.Text
		record.CertificationScheme = scheme.Selected
		record.IssueDate = parsedIssue
		record.ExpiryDate = parsedExpiry
		record.InstallLocation = installLocation.Text
		record.ResponsiblePerson = responsiblePerson.Text
		record.Notes = notes.Text
		if err := saveReferenceRecord(db, record, isEdit); err != nil {
			dialog.ShowError(err, window)
			return
		}
		window.Close()
	}

	cancel := widget.NewButton("Закрыть", func() { window.Close() })
	cancel.Importance = widget.HighImportance
	submit := widget.NewButton("Сохранить", func() { form.OnSubmit() })
	submit.Importance = widget.HighImportance
	window.SetContent(container.NewBorder(nil, container.NewHBox(cancel, submit), nil, nil, container.NewVScroll(form)))
	window.SetOnClosed(onClosed)
	window.Show()
}

func saveReferenceRecord(db *gorm.DB, record *models.SZIRecord, isEdit bool) error {
	if isEdit {
		if err := services.UpdateSziRecord(db, record.UserID, record); err != nil {
			return fmt.Errorf("не удалось сохранить запись: %v", err)
		}
		return nil
	}
	if err := services.CreateSziRecord(db, record); err != nil {
		return fmt.Errorf("не удалось создать запись: %v", err)
	}
	return nil
}

func valueOrFallback(primary, fallback string) string {
	if primary != "" {
		return primary
	}
	return fallback
}

func formatReferenceDate(primary, fallback interface{ Format(string) string }) string {
	if primary.Format("2006-01-02") != "0001-01-01" {
		return primary.Format("2006-01-02")
	}
	return fallback.Format("2006-01-02")
}
