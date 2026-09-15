package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"coursework-v2/internal/application"
	"coursework-v2/internal/domain"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func showWorkspace(myApp fyne.App, app application.App, user domain.User) {
	window := myApp.NewWindow("Реестр СЗИ | Рабочее пространство")
	window.Resize(fyne.NewSize(1100, 760))
	search := widget.NewEntry()
	search.SetPlaceHolder("Поиск по названию, сертификату, разработчику или месту...")
	status := widget.NewSelect([]string{"Все статусы", domain.StatusCurrent, domain.StatusAttention, domain.StatusExpired}, nil)
	status.SetSelected("Все статусы")
	typeFilter := widget.NewSelect(append([]string{"Все типы"}, domain.SZITypeValues...), nil)
	typeFilter.SetSelected("Все типы")
	protectionFilter := widget.NewSelect(append([]string{"Все классы СВТ"}, domain.ProtectionClassValues...), nil)
	protectionFilter.SetSelected("Все классы СВТ")
	trustFilter := widget.NewSelect(append([]string{"Все уровни доверия"}, domain.TrustLevelValues...), nil)
	trustFilter.SetSelected("Все уровни доверия")
	acFilter := widget.NewSelect(append([]string{"Все классы АС"}, domain.ACClassValues...), nil)
	acFilter.SetSelected("Все классы АС")
	schemeFilter := widget.NewSelect(append([]string{"Все схемы"}, domain.CertificationSchemeValues...), nil)
	schemeFilter.SetSelected("Все схемы")
	sortSelect := widget.NewSelect([]string{"Название: А-Я", "Срок действия", "Статус"}, nil)
	sortSelect.SetSelected("Название: А-Я")
	viewSelect := widget.NewSelect([]string{"Карточки", "Список"}, nil)
	viewSelect.SetSelected("Карточки")
	resultLabel := widget.NewLabel("")
	resultLabel.TextStyle = fyne.TextStyle{Bold: true}
	recordsView := container.NewVBox()
	scroll := container.NewVScroll(recordsView)
	var records []domain.SZIRecord
	var refresh func()
	refresh = func() {
		loaded, err := app.Records(user.ID)
		if err != nil {
			showAppError("Не удалось загрузить записи", window)
			return
		}
		records = loaded
		renderRecords(window, myApp, app, user, records, search.Text, status.Selected, typeFilter.Selected, protectionFilter.Selected, trustFilter.Selected, acFilter.Selected, schemeFilter.Selected, sortSelect.Selected, viewSelect.Selected, recordsView, resultLabel, refresh)
	}
	onFilter := func(string) {
		renderRecords(window, myApp, app, user, records, search.Text, status.Selected, typeFilter.Selected, protectionFilter.Selected, trustFilter.Selected, acFilter.Selected, schemeFilter.Selected, sortSelect.Selected, viewSelect.Selected, recordsView, resultLabel, refresh)
	}
	search.OnChanged = onFilter
	status.OnChanged = onFilter
	typeFilter.OnChanged = onFilter
	protectionFilter.OnChanged = onFilter
	trustFilter.OnChanged = onFilter
	acFilter.OnChanged = onFilter
	schemeFilter.OnChanged = onFilter
	sortSelect.OnChanged = onFilter
	viewSelect.OnChanged = onFilter
	addButton := widget.NewButtonWithIcon("Добавить СЗИ", theme.ContentAddIcon(), func() { showRecordForm(myApp, app, user.ID, nil, refresh) })
	addButton.Importance = widget.HighImportance
	refreshButton := widget.NewButtonWithIcon("Обновить", theme.ViewRefreshIcon(), refresh)
	refreshButton.Importance = widget.HighImportance
	exportButton := widget.NewButton("Экспорт CSV", func() { exportRecords(window, records) })
	exportButton.Importance = widget.HighImportance
	statisticsButton := widget.NewButton("Статистика", func() { showStatistics(myApp, app, user.ID) })
	statisticsButton.Importance = widget.HighImportance
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
	logoutButton := widget.NewButton("Выйти из учётной записи", window.Close)
	logoutButton.Importance = widget.HighImportance
	toolbar := container.NewBorder(nil, nil, nil, container.NewHBox(addButton, refreshButton, exportButton, statisticsButton), search)
	filters := container.NewVBox(container.NewHBox(status, typeFilter, protectionFilter, trustFilter), container.NewHBox(acFilter, schemeFilter, sortSelect, viewSelect, layout.NewSpacer(), clearButton))
	header := container.NewVBox(widget.NewLabelWithStyle("Рабочее пространство", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), toolbar, filters, resultLabel, widget.NewSeparator())
	window.SetContent(container.NewBorder(header, logoutButton, nil, nil, scroll))
	window.SetOnClosed(func() { showLogin(myApp, app) })
	window.Show()
	refresh()
}

func renderRecords(window fyne.Window, myApp fyne.App, app application.App, user domain.User, records []domain.SZIRecord, query, statusFilter, typeFilter, protectionFilter, trustFilter, acFilter, schemeFilter, sortMode, viewMode string, target *fyne.Container, resultLabel *widget.Label, refresh func()) {
	filtered := make([]domain.SZIRecord, 0, len(records))
	query = strings.ToLower(strings.TrimSpace(query))
	for _, record := range records {
		text := strings.ToLower(strings.Join([]string{record.Name, record.CertificateNumber, record.Vendor, record.InstallLocation, record.SZIType, record.ResponsiblePerson}, " "))
		if query != "" && !strings.Contains(text, query) {
			continue
		}
		if statusFilter != "Все статусы" && app.Status(record.ExpiryDate) != statusFilter {
			continue
		}
		if typeFilter != "Все типы" && record.SZIType != typeFilter {
			continue
		}
		if protectionFilter != "Все классы СВТ" && record.ProtectionClass != protectionFilter {
			continue
		}
		if trustFilter != "Все уровни доверия" && record.TrustLevel != trustFilter {
			continue
		}
		if acFilter != "Все классы АС" && record.ACClass != acFilter {
			continue
		}
		if schemeFilter != "Все схемы" && record.CertificationScheme != schemeFilter {
			continue
		}
		filtered = append(filtered, record)
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		if sortMode == "Срок действия" {
			return filtered[i].ExpiryDate.Before(filtered[j].ExpiryDate)
		}
		if sortMode == "Статус" {
			return app.Status(filtered[i].ExpiryDate) < app.Status(filtered[j].ExpiryDate)
		}
		return strings.ToLower(filtered[i].Name) < strings.ToLower(filtered[j].Name)
	})
	resultLabel.SetText(fmt.Sprintf("Показано записей: %d из %d", len(filtered), len(records)))
	target.Objects = nil
	for _, record := range filtered {
		if viewMode == "Список" {
			target.Add(recordRow(window, myApp, app, user, record, refresh))
		} else {
			target.Add(recordCard(window, myApp, app, user, record, refresh))
		}
	}
	if len(filtered) == 0 {
		target.Add(widget.NewLabel("По выбранным условиям записи не найдены"))
	}
	target.Refresh()
}

func recordCard(window fyne.Window, myApp fyne.App, app application.App, user domain.User, record domain.SZIRecord, refresh func()) fyne.CanvasObject {
	status := app.Status(record.ExpiryDate)
	actions := container.NewHBox(
		widget.NewButton("Изменить", func() { showRecordForm(myApp, app, user.ID, &record, refresh) }),
		widget.NewButton("Удалить", func() { confirmDelete(window, app, user.ID, record, refresh) }),
	)
	details := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("%s · Сертификат %s", record.SZIType, record.CertificateNumber)),
		widget.NewLabel("Разработчик: "+record.Vendor),
		widget.NewLabel(fmt.Sprintf("Класс СВТ: %s · Уровень доверия: %s · Класс АС: %s", valueOrDash(record.ProtectionClass), valueOrDash(record.TrustLevel), valueOrDash(record.ACClass))),
		widget.NewLabel(fmt.Sprintf("Установка: %s · Ответственный: %s", valueOrDash(record.InstallLocation), valueOrDash(record.ResponsiblePerson))),
		widget.NewLabel("Действует до: "+record.ExpiryDate.Format("02.01.2006")),
		widget.NewLabelWithStyle(status, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), actions,
	)
	return widget.NewCard(record.Name, record.SZIType, details)
}

func recordRow(window fyne.Window, myApp fyne.App, app application.App, user domain.User, record domain.SZIRecord, refresh func()) fyne.CanvasObject {
	buttons := container.NewHBox(
		widget.NewLabel(app.Status(record.ExpiryDate)),
		widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() { showRecordForm(myApp, app, user.ID, &record, refresh) }),
		widget.NewButtonWithIcon("", theme.DeleteIcon(), func() { confirmDelete(window, app, user.ID, record, refresh) }),
	)
	return container.NewBorder(nil, widget.NewSeparator(), widget.NewLabelWithStyle(record.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), buttons, widget.NewLabel(fmt.Sprintf("%s · %s · до %s", record.SZIType, record.Vendor, record.ExpiryDate.Format("02.01.2006"))))
}
func confirmDelete(window fyne.Window, app application.App, userID uint, record domain.SZIRecord, refresh func()) {
	dialog.ShowConfirm("Удаление записи", "Удалить «"+record.Name+"»?", func(ok bool) {
		if ok {
			if err := app.Delete(userID, record.ID); err != nil {
				showAppError(err.Error(), window)
				return
			}
			refresh()
		}
	}, window)
}
func valueOrDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "не указано"
	}
	return value
}

func referenceType(record domain.SZIRecord) string            { return record.SZIType }
func referenceCertificate(record domain.SZIRecord) string     { return record.CertificateNumber }
func referenceVendor(record domain.SZIRecord) string          { return record.Vendor }
func referenceLocation(record domain.SZIRecord) string        { return record.InstallLocation }
func referenceProtectionClass(record domain.SZIRecord) string { return record.ProtectionClass }
func referenceExpiryDate(record domain.SZIRecord) time.Time   { return record.ExpiryDate }
