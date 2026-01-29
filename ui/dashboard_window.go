package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
	"szi-registry/models"
	"szi-registry/services"
)

// ShowDashboardWindow отображает главное окно с панелью статистики
func ShowDashboardWindow(myApp fyne.App, username string, db *gorm.DB) {
	myWindow := myApp.NewWindow("Реестр СЗИ от НСД — Главная")
	myWindow.Resize(fyne.NewSize(1000, 700))

	// Получаем ID пользователя по имени
	user, err := services.GetUserByUsername(db, username)
	if err != nil || user == nil {
		// Если не удалось получить пользователя, используем 0
		user = &models.User{ID: 0}
	}

	// Получаем статистику
	stats, err := services.GetSZIStatistics(db, uint(user.ID))
	if err != nil {
		dialog := widget.NewModalPopUp(
			widget.NewLabel("Ошибка при получении статистики: "+err.Error()),
			myWindow.Canvas(),
		)
		dialog.Show()
		return
	}

	// Создаем контейнер для статистики
	statContainer := container.NewVBox()

	// Добавляем общую информацию
	totalLabel := widget.NewLabel(fmt.Sprintf("Всего СЗИ: %d", stats.TotalCount))
	totalLabel.TextStyle.Bold = true
	statContainer.Add(totalLabel)

	// Добавляем информацию о сертификатах
	certHeader := widget.NewLabel("Сертификаты:")
	certHeader.TextStyle.Bold = true
	statContainer.Add(certHeader)

	activeCertLabel := widget.NewLabel(fmt.Sprintf("• Действующие: %d", stats.CertificateStats.Active))
	expiredCertLabel := widget.NewLabel(fmt.Sprintf("• Истёкшие: %d", stats.CertificateStats.Expired))
	noCertLabel := widget.NewLabel(fmt.Sprintf("• Без сертификата: %d", stats.CertificateStats.NoCert))
	statContainer.Add(activeCertLabel)
	statContainer.Add(expiredCertLabel)
	statContainer.Add(noCertLabel)

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

	// Кнопка перехода к реестру СЗИ
	gotoRegistryBtn := widget.NewButton("Перейти к реестру СЗИ", func() {
		ShowSZIRegistryWindow(myApp, username, db)
		myWindow.Close()
	})

	// Кнопка выхода
	logoutBtn := widget.NewButton("Выход", func() {
		myWindow.Close()
		ShowLoginWindow(myApp, db)
	})

	// Контейнер для кнопок
	buttonsContainer := container.NewHBox(gotoRegistryBtn, logoutBtn)

	// Основной контейнер
	mainContainer := container.NewBorder(
		nil, // верх
		buttonsContainer, // низ
		nil, // лево
		nil, // право
		container.NewPadded(statContainer), // центр
	)

	myWindow.SetContent(mainContainer)
	myWindow.Show()
}