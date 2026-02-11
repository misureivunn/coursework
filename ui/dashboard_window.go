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

//отображает главное окно с панелью статистики
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

	statContainer := container.NewVBox()

	totalLabel := widget.NewLabel(fmt.Sprintf("Всего СЗИ: %d", stats.TotalCount))
	totalLabel.TextStyle.Bold = true
	statContainer.Add(totalLabel)

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

	//переход к реестру СЗИ
	gotoRegistryBtn := widget.NewButton("Перейти к реестру СЗИ", func() {
		ShowSZIRegistryWindow(myApp, username, db)
		myWindow.Close()
	})

	logoutBtn := widget.NewButton("Выход", func() {
		myWindow.Close()
		ShowLoginWindow(myApp, db)
	})

	buttonsContainer := container.NewHBox(gotoRegistryBtn, logoutBtn)

	mainContainer := container.NewBorder(
		nil,
		buttonsContainer, 
		nil, 
		nil, 
		container.NewPadded(statContainer),
	)

	myWindow.SetContent(mainContainer)
	myWindow.Show()
}