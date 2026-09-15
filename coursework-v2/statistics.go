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
	"fyne.io/fyne/v2/widget"
)

func showStatistics(myApp fyne.App, app application.App, userID uint) {
	window := myApp.NewWindow("Статистика реестра")
	window.Resize(fyne.NewSize(620, 620))
	records, err := app.Records(userID)
	if err != nil {
		dialog.ShowError(fmt.Errorf("не удалось загрузить статистику: %v", err), window)
		return
	}

	statusCounts := map[string]int{}
	typeCounts := map[string]int{}
	classCounts := map[string]int{}
	trustCounts := map[string]int{}
	for _, record := range records {
		statusCounts[domain.StatusFor(record.ExpiryDate, time.Now())]++
		if record.SZIType != "" {
			typeCounts[record.SZIType]++
		}
		if record.ProtectionClass != "" {
			classCounts[record.ProtectionClass]++
		}
		if record.TrustLevel != "" {
			trustCounts[record.TrustLevel]++
		}
	}
	content := container.NewVBox(
		widget.NewLabelWithStyle(fmt.Sprintf("Всего записей: %d", len(records)), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(fmt.Sprintf("Актуально: %d · Требует внимания: %d · Просрочено: %d", statusCounts[domain.StatusCurrent], statusCounts[domain.StatusAttention], statusCounts[domain.StatusExpired])),
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
