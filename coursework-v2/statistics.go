package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"szi-registry/models"
	"szi-registry/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
)

func showStatistics(myApp fyne.App, db *gorm.DB, userID uint) {
	window := myApp.NewWindow("Статистика реестра")
	window.Resize(fyne.NewSize(620, 620))
	var records []models.SZIRecord
	if err := db.Where("user_id = ?", userID).Find(&records).Error; err != nil {
		dialog.ShowError(fmt.Errorf("не удалось загрузить статистику: %v", err), window)
		return
	}

	statusCounts := map[string]int{}
	typeCounts := map[string]int{}
	classCounts := map[string]int{}
	trustCounts := map[string]int{}
	for _, record := range records {
		statusCounts[services.CalculateSziStatus(referenceExpiryDate(record), time.Now())]++
		if referenceType(record) != "" {
			typeCounts[referenceType(record)]++
		}
		if referenceProtectionClass(record) != "" {
			classCounts[referenceProtectionClass(record)]++
		}
		if record.TrustLevel != "" {
			trustCounts[record.TrustLevel]++
		}
	}
	content := container.NewVBox(
		widget.NewLabelWithStyle(fmt.Sprintf("Всего записей: %d", len(records)), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(fmt.Sprintf("Актуально: %d · Требует внимания: %d · Просрочено: %d", statusCounts[services.StatusCurrent], statusCounts[services.StatusAttention], statusCounts[services.StatusExpired])),
		widget.NewSeparator(),
		statSection("Типы СЗИ", typeCounts),
		statSection("Классы СВТ", classCounts),
		statSection("Уровни доверия", trustCounts),
	)
	closeButton := widget.NewButton("Закрыть", window.Close)
	closeButton.Importance = widget.HighImportance
	window.SetContent(container.NewBorder(nil, closeButton, nil, nil, container.NewVScroll(content)))
	window.Show()
}

func statSection(title string, values map[string]int) fyne.CanvasObject {
	rows := []string{}
	for key, value := range values {
		rows = append(rows, fmt.Sprintf("%s: %d", key, value))
	}
	sort.Strings(rows)
	if len(rows) == 0 {
		rows = append(rows, "Нет данных")
	}
	return container.NewVBox(widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), widget.NewLabel(strings.Join(rows, "\n")))
}

func referenceIssueDate(record models.SZIRecord) time.Time {
	if !record.IssueDate.IsZero() {
		return record.IssueDate
	}
	return record.CertIssueDate
}

func now() time.Time { return time.Now() }
