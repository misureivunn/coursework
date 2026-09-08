package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"szi-registry/models"
	"szi-registry/services"
	legacyui "szi-registry/ui"
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
	search.SetPlaceHolder("Поиск по названию, сертификату, производителю или месту...")
	status := widget.NewSelect([]string{"Все статусы", services.StatusCurrent, services.StatusAttention, services.StatusExpired}, nil)
	status.SetSelected("Все статусы")
	typeFilter := widget.NewSelect([]string{"Все типы", "СКЗИ", "ОС", "СЗИ КС", "СЗИ СКЗИ", "Другое"}, nil)
	typeFilter.SetSelected("Все типы")
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
		renderRecords(window, myApp, db, user, records, search.Text, status.Selected, typeFilter.Selected, sortSelect.Selected, viewSelect.Selected, recordsView, resultLabel, refresh)
	}

	search.OnChanged = func(string) {
		renderRecords(window, myApp, db, user, records, search.Text, status.Selected, typeFilter.Selected, sortSelect.Selected, viewSelect.Selected, recordsView, resultLabel, refresh)
	}
	status.OnChanged = func(string) {
		renderRecords(window, myApp, db, user, records, search.Text, status.Selected, typeFilter.Selected, sortSelect.Selected, viewSelect.Selected, recordsView, resultLabel, refresh)
	}
	typeFilter.OnChanged = func(string) {
		renderRecords(window, myApp, db, user, records, search.Text, status.Selected, typeFilter.Selected, sortSelect.Selected, viewSelect.Selected, recordsView, resultLabel, refresh)
	}
	sortSelect.OnChanged = func(string) {
		renderRecords(window, myApp, db, user, records, search.Text, status.Selected, typeFilter.Selected, sortSelect.Selected, viewSelect.Selected, recordsView, resultLabel, refresh)
	}
	viewSelect.OnChanged = func(string) {
		renderRecords(window, myApp, db, user, records, search.Text, status.Selected, typeFilter.Selected, sortSelect.Selected, viewSelect.Selected, recordsView, resultLabel, refresh)
	}

	addButton := widget.NewButtonWithIcon("Добавить СЗИ", theme.ContentAddIcon(), func() {
		legacyui.ShowAddSziWindow(myApp, user.ID, db)
	})
	addButton.Importance = widget.HighImportance
	refreshButton := widget.NewButtonWithIcon("Обновить", theme.ViewRefreshIcon(), refresh)
	clearButton := widget.NewButton("Сбросить", func() {
		search.SetText("")
		status.SetSelected("Все статусы")
		typeFilter.SetSelected("Все типы")
		sortSelect.SetSelected("Название: А-Я")
		refresh()
	})

	toolbar := container.NewBorder(nil, nil, nil, container.NewHBox(addButton, refreshButton), search)
	filters := container.NewHBox(status, typeFilter, sortSelect, viewSelect, layout.NewSpacer(), clearButton)
	header := container.NewVBox(
		widget.NewLabelWithStyle("Рабочее пространство", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Управляйте реестром без лишних переходов"),
		toolbar,
		filters,
		resultLabel,
		widget.NewSeparator(),
	)

	window.SetContent(container.NewBorder(header, nil, nil, nil, scroll))
	window.SetOnClosed(func() { showLogin(myApp, db) })
	window.Show()
	refresh()
}

func renderRecords(window fyne.Window, myApp fyne.App, db *gorm.DB, user *models.User, records []models.SZIRecord, query, statusFilter, typeFilter, sortMode, viewMode string, target *fyne.Container, resultLabel *widget.Label, refresh func()) {
	filtered := make([]models.SZIRecord, 0, len(records))
	query = strings.ToLower(strings.TrimSpace(query))
	for _, record := range records {
		currentStatus := services.CalculateSziStatus(record.CertExpiryDate, time.Now())
		searchText := strings.ToLower(strings.Join([]string{record.Name, record.CertNumber, record.Manufacturer, record.Location, record.Type}, " "))
		if query != "" && !strings.Contains(searchText, query) {
			continue
		}
		if statusFilter != "" && statusFilter != "Все статусы" && currentStatus != statusFilter {
			continue
		}
		if typeFilter != "" && typeFilter != "Все типы" && record.Type != typeFilter {
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
	meta := fmt.Sprintf("%s  ·  Сертификат %s", record.Type, record.CertNumber)
	details := container.NewVBox(
		widget.NewLabel(meta),
		widget.NewLabel(fmt.Sprintf("Производитель: %s", valueOrDash(record.Manufacturer))),
		widget.NewLabel(fmt.Sprintf("Установка: %s", valueOrDash(record.Location))),
		widget.NewLabel(fmt.Sprintf("Действует до: %s", record.CertExpiryDate.Format("02.01.2006"))),
		widget.NewLabelWithStyle(status, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(
			widget.NewButtonWithIcon("Изменить", theme.DocumentCreateIcon(), func() { legacyui.ShowEditSziWindow(myApp, record, db) }),
			widget.NewButtonWithIcon("Удалить", theme.DeleteIcon(), func() { confirmRecordDelete(window, db, user, record, refresh) }),
		),
	)
	return widget.NewCard(record.Name, meta, details)
}

func newRecordRow(window fyne.Window, myApp fyne.App, db *gorm.DB, user *models.User, record models.SZIRecord, refresh func()) fyne.CanvasObject {
	status := services.CalculateSziStatus(record.CertExpiryDate, time.Now())
	edit := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() { legacyui.ShowEditSziWindow(myApp, record, db) })
	delete := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() { confirmRecordDelete(window, db, user, record, refresh) })
	return container.NewBorder(nil, widget.NewSeparator(), widget.NewLabelWithStyle(record.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), container.NewHBox(widget.NewLabel(status), edit, delete), widget.NewLabel(fmt.Sprintf("%s · %s · до %s", record.Type, valueOrDash(record.Manufacturer), record.CertExpiryDate.Format("02.01.2006"))))
}

func confirmRecordDelete(window fyne.Window, db *gorm.DB, user *models.User, record models.SZIRecord, refresh func()) {
	dialog.ShowConfirm("Удаление записи", "Удалить «"+record.Name+"»?", func(ok bool) {
		if !ok {
			return
		}
		if err := services.DeleteSziRecord(db, user.ID, record.ID); err != nil {
			appui.ShowErrorDialog(err.Error(), window.Canvas())
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
