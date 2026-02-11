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

// ShowNotificationsWindow отображает окно с уведомлениями
func ShowNotificationsWindow(myApp fyne.App, username string, db *gorm.DB) {
	myWindow := myApp.NewWindow("🔔 Уведомления")
	myWindow.Resize(fyne.NewSize(800, 600))
	
	user, _ := services.GetUserByUsername(db, username)
	if user == nil {
		return
	}
	
	// Получаем уведомления
	notifications, err := services.GetUserNotifications(db, user.ID, false)
	if err != nil {
		notifications = []models.Notification{}
	}
	
	var content *fyne.Container
	
	if len(notifications) == 0 {
		// Нет уведомлений
		emptyLabel := widget.NewLabel("📭 Уведомлений пока нет")
		emptyLabel.Alignment = fyne.TextAlignCenter
		emptyLabel.TextStyle.Bold = true
		
		content = container.NewCenter(emptyLabel)
	} else {
		// Создаём список уведомлений
		notifList := container.NewVBox()
		
		for _, notif := range notifications {
			notifCard := createNotificationCard(notif, db, myWindow, myApp, username)
			notifList.Add(notifCard)
			notifList.Add(widget.NewSeparator())
		}
		
		content = container.NewVScroll(notifList)
	}
	
	// Кнопка "Отметить все как прочитанные"
	markAllBtn := widget.NewButton("✅ Отметить все как прочитанные", func() {
		services.MarkAllNotificationsAsRead(db, user.ID)
		myWindow.Close()
		ShowNotificationsWindow(myApp, username, db)
	})
	
	backBtn := widget.NewButton("◀ Назад", func() {
		myWindow.Close()
		ShowDashboardWindow(myApp, username, db)
	})
	
	toolbar := container.NewHBox(backBtn, markAllBtn)
	
	finalContent := container.NewBorder(toolbar, nil, nil, nil, content)
	myWindow.SetContent(finalContent)
	myWindow.Show()
}

func createNotificationCard(notif models.Notification, db *gorm.DB, window fyne.Window, myApp fyne.App, username string) *fyne.Container {
	titleLabel := widget.NewLabel(notif.Title)
	titleLabel.TextStyle.Bold = true
	
	// Устанавливаем важность в зависимости от приоритета
	if notif.Priority == 4 {
		titleLabel.Importance = widget.DangerImportance
	} else if notif.Priority == 3 {
		titleLabel.Importance = widget.WarningImportance
	}
	
	messageLabel := widget.NewLabel(notif.Message)
	messageLabel.Wrapping = fyne.TextWrapWord
	
	timeLabel := widget.NewLabel(notif.CreatedAt.Format("2006-01-02 15:04"))
	timeLabel.TextStyle.Italic = true
	
	markReadBtn := widget.NewButton("✓ Прочитано", func() {
		services.MarkNotificationAsRead(db, notif.ID)
		window.Close()
		ShowNotificationsWindow(myApp, username, db)
	})
	
	deleteBtn := widget.NewButton("🗑", func() {
		services.DeleteNotification(db, notif.ID)
		window.Close()
		ShowNotificationsWindow(myApp, username, db)
	})
	
	buttons := container.NewHBox(markReadBtn, deleteBtn)
	
	// Фон для непрочитанных
	card := container.NewVBox(
		titleLabel,
		messageLabel,
		timeLabel,
		buttons,
	)
	
	if !notif.IsRead {
		unreadLabel := widget.NewLabel("● Непрочитано")
		unreadLabel.Importance = widget.HighImportance
		card.Add(unreadLabel)
	}
	
	return card
}
