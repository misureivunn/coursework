package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"szi-registry/models"
	"szi-registry/services"
	appui "szi-registry/utils/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

func showWorkspace(myApp fyne.App, db *gorm.DB, user *models.User) {
	window := myApp.NewWindow("Реестр СЗИ | Рабочее пространство")
	window.Resize(fyne.NewSize(1100, 760))

	search := widget.NewEntry()
	search.SetPlaceHolder("Поиск по названию, сертификату, разработчику или месту...")
	status := widget.NewSelect([]string{"Все статусы", services.StatusCurrent, services.StatusAttention, services.StatusExpired}, nil)
	status.SetSelected("Все статусы")
	typeFilter := widget.NewSelect(append([]string{"Все типы"}, models.SZITypeValues...), nil)
	typeFilter.SetSelected("Все типы")
	protectionFilter := widget.NewSelect(append([]string{"Все классы СВТ"}, models.ProtectionClassValues...), nil)
	protectionFilter.SetSelected("Все классы СВТ")
	trustFilter := widget.NewSelect(append([]string{"Все уровни доверия"}, models.TrustLevelValues...), nil)
	trustFilter.SetSelected("Все уровни доверия")
	acFilter := widget.NewSelect(append([]string{"Все классы АС"}, models.ACClassValues...), nil)
	acFilter.SetSelected("Все классы АС")
	schemeFilter := widget.NewSelect(append([]string{"Все схемы"}, models.CertificationSchemeValues...), nil)
	schemeFilter.SetSelected("Все схемы")
	sortSelect := widget.NewSelect([]string{"Название: А-Я", "Срок действия", "Статус"}, nil)
	sortSelect.SetSelected("Название: А-Я")
	viewSelect := widget.NewSelect([]string{"Карточки", "Список"}, nil)
	viewSelect.SetSelected("Карточки")

	resultLabel := widget.NewLabel("")
	resultLabel.TextStyle = fyne.TextStyle{Bold: true}
	recordsView := container.NewVBox()
	scroll := container.NewVScroll(recordsView)

	var records []models.SZIRecord
	var refresh func()
	refresh = func() {
		loaded, err := services.GetUserRecords(db, user.ID)
		if err != nil {
			appui.ShowErrorDialog("Не удалось загрузить записи: "+err.Error(), window.Canvas())
			return
		}
		records = loaded
		renderRecords(window, myApp, db, user, records, search.Text, status.Selected, typeFilter.Selected, protectionFilter.Selected, trustFilter.Selected, acFilter.Selected, schemeFilter.Selected, sortSelect.Selected, viewSelect.Selected, recordsView, resultLabel, refresh)
	}

	search.OnChanged = func(string) {
		renderRecords(window, myApp, db, user, records, search.Text, status.Selected, typeFilter.Selected, protectionFilter.Selected, trustFilter.Selected, acFilter.Selected, schemeFilter.Selected, sortSelect.Selected, viewSelect.Selected, recordsView, resultLabel, refresh)
	}
	status.OnChanged = func(string) {
		renderRecords(window, myApp, db, user, records, search.Text, status.Selected, typeFilter.Selected, protectionFilter.Selected, trustFilter.Selected, acFilter.Selected, schemeFilter.Selected, sortSelect.Selected, viewSelect.Selected, recordsView, resultLabel, refresh)
	}
	typeFilter.OnChanged = func(string) {
		renderRecords(window, myApp, db, user, records, search.Text, status.Selected, typeFilter.Selected, protectionFilter.Selected, trustFilter.Selected, acFilter.Selected, schemeFilter.Selected, sortSelect.Selected, viewSelect.Selected, recordsView, resultLabel, refresh)
	}
	sortSelect.OnChanged = func(string) {
		renderRecords(window, myApp, db, user, records, search.Text, status.Selected, typeFilter.Selected, protectionFilter.Selected, trustFilter.Selected, acFilter.Selected, schemeFilter.Selected, sortSelect.Selected, viewSelect.Selected, recordsView, resultLabel, refresh)
	}
	viewSelect.OnChanged = func(string) {
		renderRecords(window, myApp, db, user, records, search.Text, status.Selected, typeFilter.Selected, protectionFilter.Selected, trustFilter.Selected, acFilter.Selected, schemeFilter.Selected, sortSelect.Selected, viewSelect.Selected, recordsView, resultLabel, refresh)
	}
	for _, filter := range []*widget.Select{protectionFilter, trustFilter, acFilter, schemeFilter} {
		filter.OnChanged = func(string) {
			renderRecords(window, myApp, db, user, records, search.Text, status.Selected, typeFilter.Selected, protectionFilter.Selected, trustFilter.Selected, acFilter.Selected, schemeFilter.Selected, sortSelect.Selected, viewSelect.Selected, recordsView, resultLabel, refresh)
		}
	}

	addButton := widget.NewButtonWithIcon("Добавить СЗИ", theme.ContentAddIcon(), func() {
		showRecordForm(myApp, db, user.ID, nil, refresh)
	})
	addButton.Importance = widget.HighImportance
	refreshButton := widget.NewButtonWithIcon("Обновить", theme.ViewRefreshIcon(), refresh)
	refreshButton.Importance = widget.HighImportance
	clearButton := widget.NewButton("Сбросить", func() {
		search.SetText("")
		status.SetSelected("Все статусы")
		typeFilter.SetSelected("Все типы")
		protectionFilter.SetSelected("Все классы СВТ")
		trustFilter.SetSelected("Все уровни доверия")
		acFilter.SetSelected("Все классы АС")
		schemeFilter.SetSelected("Все схемы")
		sortSelect.SetSelected("Название: А-Я")
		refresh()
	})
	clearButton.Importance = widget.HighImportance
	logoutButton := widget.NewButton("Выйти из учётной записи", func() { window.Close() })
	logoutButton.Importance = widget.HighImportance

	toolbar := container.NewBorder(nil, nil, nil, container.NewHBox(addButton, refreshButton), search)
	filters := container.NewVBox(
		container.NewHBox(status, typeFilter, protectionFilter, trustFilter),
		container.NewHBox(acFilter, schemeFilter, sortSelect, viewSelect, layout.NewSpacer(), clearButton),
	)
	header := container.NewVBox(
		widget.NewLabelWithStyle("Рабочее пространство", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Управляйте реестром без лишних переходов"),
		toolbar,
		filters,
		resultLabel,
		widget.NewSeparator(),
	)

	window.SetContent(container.NewBorder(header, logoutButton, nil, nil, scroll))
	window.SetOnClosed(func() { showLogin(myApp, db) })
	window.Show()
	refresh()
}

func renderRecords(window fyne.Window, myApp fyne.App, db *gorm.DB, user *models.User, records []models.SZIRecord, query, statusFilter, typeFilter, protectionFilter, trustFilter, acFilter, schemeFilter, sortMode, viewMode string, target *fyne.Container, resultLabel *widget.Label, refresh func()) {
	filtered := make([]models.SZIRecord, 0, len(records))
	query = strings.ToLower(strings.TrimSpace(query))
	for _, record := range records {
		currentStatus := services.CalculateSziStatus(record.CertExpiryDate, time.Now())
		searchText := strings.ToLower(strings.Join([]string{record.Name, referenceCertificate(record), referenceVendor(record), referenceLocation(record), referenceType(record), record.ResponsiblePerson}, " "))
		if query != "" && !strings.Contains(searchText, query) {
			continue
		}
		if statusFilter != "" && statusFilter != "Все статусы" && currentStatus != statusFilter {
			continue
		}
		if typeFilter != "" && typeFilter != "Все типы" && referenceType(record) != typeFilter {
			continue
		}
		if protectionFilter != "" && protectionFilter != "Все классы СВТ" && referenceProtectionClass(record) != protectionFilter {
			continue
		}
		if trustFilter != "" && trustFilter != "Все уровни доверия" && record.TrustLevel != trustFilter {
			continue
		}
		if acFilter != "" && acFilter != "Все классы АС" && record.ACClass != acFilter {
			continue
		}
		if schemeFilter != "" && schemeFilter != "Все схемы" && record.CertificationScheme != schemeFilter {
			continue
		}
		filtered = append(filtered, record)
	}

	sort.SliceStable(filtered, func(i, j int) bool {
		switch sortMode {
		case "Срок действия":
			return filtered[i].CertExpiryDate.Before(filtered[j].CertExpiryDate)
		case "Статус":
			return services.CalculateSziStatus(filtered[i].CertExpiryDate, time.Now()) < services.CalculateSziStatus(filtered[j].CertExpiryDate, time.Now())
		default:
			return strings.ToLower(filtered[i].Name) < strings.ToLower(filtered[j].Name)
		}
	})

	resultLabel.SetText(fmt.Sprintf("Показано записей: %d из %d", len(filtered), len(records)))
	target.Objects = nil
	if len(filtered) == 0 {
		target.Add(widget.NewLabel("По выбранным условиям записи не найдены"))
	} else {
		for _, record := range filtered {
			if viewMode == "Список" {
				target.Add(newRecordRow(window, myApp, db, user, record, refresh))
			} else {
				target.Add(newRecordCard(window, myApp, db, user, record, refresh))
			}
		}
	}
	target.Refresh()
}

func newRecordCard(window fyne.Window, myApp fyne.App, db *gorm.DB, user *models.User, record models.SZIRecord, refresh func()) fyne.CanvasObject {
	status := services.CalculateSziStatus(record.CertExpiryDate, time.Now())
	meta := fmt.Sprintf("%s  ·  Сертификат %s", referenceType(record), referenceCertificate(record))
	details := container.NewVBox(
		widget.NewLabel(meta),
		widget.NewLabel(fmt.Sprintf("Разработчик / заявитель: %s", valueOrDash(referenceVendor(record)))),
		widget.NewLabel(fmt.Sprintf("Класс СВТ: %s · Уровень доверия: %s · Класс АС: %s", valueOrDash(referenceProtectionClass(record)), valueOrDash(record.TrustLevel), valueOrDash(record.ACClass))),
		widget.NewLabel(fmt.Sprintf("Установка: %s · Ответственный: %s", valueOrDash(referenceLocation(record)), valueOrDash(record.ResponsiblePerson))),
		widget.NewLabel(fmt.Sprintf("Действует до: %s", referenceExpiryDate(record).Format("02.01.2006"))),
		widget.NewLabelWithStyle(status, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(
			buttonWithImportance(widget.NewButtonWithIcon("Изменить", theme.DocumentCreateIcon(), func() { showRecordForm(myApp, db, user.ID, &record, refresh) })),
			buttonWithImportance(widget.NewButtonWithIcon("Удалить", theme.DeleteIcon(), func() { confirmRecordDelete(window, db, user, record, refresh) })),
		),
	)
	return widget.NewCard(record.Name, meta, details)
}

func newRecordRow(window fyne.Window, myApp fyne.App, db *gorm.DB, user *models.User, record models.SZIRecord, refresh func()) fyne.CanvasObject {
	status := services.CalculateSziStatus(record.CertExpiryDate, time.Now())
	edit := buttonWithImportance(widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() { showRecordForm(myApp, db, user.ID, &record, refresh) }))
	delete := buttonWithImportance(widget.NewButtonWithIcon("", theme.DeleteIcon(), func() { confirmRecordDelete(window, db, user, record, refresh) }))
	return container.NewBorder(nil, widget.NewSeparator(), widget.NewLabelWithStyle(record.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), container.NewHBox(widget.NewLabel(status), edit, delete), widget.NewLabel(fmt.Sprintf("%s · %s · до %s", referenceType(record), valueOrDash(referenceVendor(record)), referenceExpiryDate(record).Format("02.01.2006"))))
}

func referenceType(record models.SZIRecord) string {
	return valueOrFallback(record.SZIType, record.Type)
}
func referenceCertificate(record models.SZIRecord) string {
	return valueOrFallback(record.CertificateNumber, record.CertNumber)
}
func referenceVendor(record models.SZIRecord) string {
	return valueOrFallback(record.Vendor, record.Manufacturer)
}
func referenceLocation(record models.SZIRecord) string {
	return valueOrFallback(record.InstallLocation, record.Location)
}
func referenceProtectionClass(record models.SZIRecord) string {
	return valueOrFallback(record.ProtectionClass, record.ClassProtection)
}
func referenceExpiryDate(record models.SZIRecord) time.Time {
	if !record.ExpiryDate.IsZero() {
		return record.ExpiryDate
	}
	return record.CertExpiryDate
}

func confirmRecordDelete(window fyne.Window, db *gorm.DB, user *models.User, record models.SZIRecord, refresh func()) {
	dialog.ShowConfirm("Удаление записи", "Удалить «"+record.Name+"»?", func(ok bool) {
		if !ok {
			return
		}
		if err := services.DeleteSziRecord(db, user.ID, record.ID); err != nil {
			appui.ShowErrorDialog("Не удалось удалить запись СЗИ: "+err.Error(), window.Canvas())
			return
		}
		refresh()
	}, window)
}

func valueOrDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "не указано"
	}
	return value
}

func buttonWithImportance(button *widget.Button) *widget.Button {
	button.Importance = widget.HighImportance
	return button
}
