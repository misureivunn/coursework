package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
	"szi-registry/models"
	"szi-registry/services"
	"szi-registry/utils/ui"
)

// Отображает окно со статистикой по СЗИ
func showStatisticsWindow(myApp fyne.App, userID uint, db *gorm.DB) {
	window := myApp.NewWindow("Статистика СЗИ")
	window.Resize(fyne.NewSize(600, 500))

	// Получаем статистику
	stats, err := services.GetSZIStatistics(db, userID)
	if err != nil {
		dialog := widget.NewModalPopUp(
			widget.NewLabel("Ошибка при получении статистики: "+err.Error()),
			window.Canvas(),
		)
		dialog.Show()
		return
	}

	// Создаем контейнер для статистики
	statContainer := container.NewVBox()

	// Добавляем общую информацию
	totalLabel := widget.NewLabel(fmt.Sprintf("Всего СЗИ: %d", stats.TotalCount))
	statContainer.Add(totalLabel)

	// Добавляем информацию о сертификатах
	certHeader := widget.NewLabel("Сертификаты:")
	certHeader.TextStyle.Bold = true
	statContainer.Add(certHeader)

	// Получаем статистику по сертификатам из всех записей
	allRecords, err := services.GetUserRecords(db, userID)
	if err != nil {
		allRecords = []models.SZIRecord{} // Используем пустой массив в случае ошибки
	}
	active, expired, noCert := ui.GetCertificateStats(allRecords)

	activeCertLabel := widget.NewLabel(fmt.Sprintf("• Действующие: %d", active))
	expiredCertLabel := widget.NewLabel(fmt.Sprintf("• Истёкшие: %d", expired))
	noCertLabel := widget.NewLabel(fmt.Sprintf("• Без сертификата: %d", noCert))
	statContainer.Add(activeCertLabel)
	statContainer.Add(expiredCertLabel)
	statContainer.Add(noCertLabel)

	// Добавляем информацию по классам защиты
	classHeader := widget.NewLabel("Классы защиты:")
	classHeader.TextStyle.Bold = true
	statContainer.Add(classHeader)

	for protectionClass, count := range stats.ProtectionClassStats {
		if protectionClass != "" {
			classLabel := widget.NewLabel(fmt.Sprintf("• %s: %d", protectionClass, count))
			statContainer.Add(classLabel)
		}
	}

	// Добавляем информацию по типам СЗИ
	typeHeader := widget.NewLabel("Типы СЗИ:")
	typeHeader.TextStyle.Bold = true
	statContainer.Add(typeHeader)

	for szType, count := range stats.TypeStats {
		if szType != "" {
			typeLabel := widget.NewLabel(fmt.Sprintf("• %s: %d", szType, count))
			statContainer.Add(typeLabel)
		}
	}

	// Добавляем информацию по назначению
	purposeHeader := widget.NewLabel("Назначение:")
	purposeHeader.TextStyle.Bold = true
	statContainer.Add(purposeHeader)

	for purpose, count := range stats.PurposeStats {
		if purpose != "" {
			purposeLabel := widget.NewLabel(fmt.Sprintf("• %s: %d", purpose, count))
			statContainer.Add(purposeLabel)
		}
	}

	// Добавляем информацию по типу развертывания
	deploymentHeader := widget.NewLabel("Тип развертывания:")
	deploymentHeader.TextStyle.Bold = true
	statContainer.Add(deploymentHeader)

	for deploymentType, count := range stats.DeploymentTypeStats {
		if deploymentType != "" {
			deploymentLabel := widget.NewLabel(fmt.Sprintf("• %s: %d", deploymentType, count))
			statContainer.Add(deploymentLabel)
		}
	}

	// Кнопка закрытия
	closeBtn := widget.NewButton("Закрыть", func() {
		window.Close()
	})

	scrollContainer := container.NewVScroll(statContainer)
	scrollContainer.SetMinSize(fyne.NewSize(580, 400))

	content := container.NewBorder(nil, closeBtn, nil, nil, scrollContainer)
	window.SetContent(content)
	window.Show()
}