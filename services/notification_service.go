package services

import (
	"fmt"
	"gorm.io/gorm"
	"szi-registry/models"
	"szi-registry/repositories"
	"time"
)

// GenerateExpiryNotifications создаёт уведомления об истекающих сертификатах
func GenerateExpiryNotifications(db *gorm.DB) error {
	// Получаем все записи СЗИ
	var records []models.SZIRecord
	if err := db.Find(&records).Error; err != nil {
		return err
	}

	repo := repositories.NewNotificationRepository(db)
	now := time.Now()

	for _, record := range records {
		daysUntilExpiry := int(record.CertExpiryDate.Sub(now).Hours() / 24)

		// Проверяем, есть ли уже уведомление для этой записи
		var existingNotification models.Notification
		result := db.Where("record_id = ? AND is_read = ?", record.ID, false).First(&existingNotification)

		if result.Error == nil {
			// Уведомление уже существует
			continue
		}

		var notification *models.Notification

		if daysUntilExpiry < 0 {
			// Сертификат просрочен
			notification = &models.Notification{
				UserID:   record.UserID,
				RecordID: record.ID,
				Type:     "EXPIRED",
				Title:    "❌ Сертификат истёк",
				Message:  fmt.Sprintf("Сертификат СЗИ '%s' истёк %d дней назад. Использование запрещено ФСТЭК!", record.Name, -daysUntilExpiry),
				Priority: 4,
				IsRead:   false,
			}
		} else if daysUntilExpiry <= 7 {
			// Критично (< 7 дней)
			notification = &models.Notification{
				UserID:   record.UserID,
				RecordID: record.ID,
				Type:     "CRITICAL",
				Title:    "🚨 Критическое предупреждение",
				Message:  fmt.Sprintf("До истечения сертификата СЗИ '%s' осталось всего %d дней!", record.Name, daysUntilExpiry),
				Priority: 3,
				IsRead:   false,
			}
		} else if daysUntilExpiry <= 30 {
			// Скоро истекает (7-30 дней)
			notification = &models.Notification{
				UserID:   record.UserID,
				RecordID: record.ID,
				Type:     "EXPIRING_SOON",
				Title:    "⚠️ Скоро истечёт сертификат",
				Message:  fmt.Sprintf("До истечения сертификата СЗИ '%s' осталось %d дней. Требуется продление.", record.Name, daysUntilExpiry),
				Priority: 2,
				IsRead:   false,
			}
		}

		if notification != nil {
			if err := repo.Create(notification); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetUserNotifications возвращает уведомления пользователя
func GetUserNotifications(db *gorm.DB, userID uint, unreadOnly bool) ([]models.Notification, error) {
	repo := repositories.NewNotificationRepository(db)
	return repo.FindByUserID(userID, unreadOnly)
}

// MarkNotificationAsRead помечает уведомление как прочитанное
func MarkNotificationAsRead(db *gorm.DB, notificationID uint) error {
	repo := repositories.NewNotificationRepository(db)
	return repo.MarkAsRead(notificationID)
}

// MarkAllNotificationsAsRead помечает все уведомления пользователя как прочитанные
func MarkAllNotificationsAsRead(db *gorm.DB, userID uint) error {
	repo := repositories.NewNotificationRepository(db)
	return repo.MarkAllAsRead(userID)
}

// DeleteNotification удаляет уведомление
func DeleteNotification(db *gorm.DB, notificationID uint) error {
	repo := repositories.NewNotificationRepository(db)
	return repo.Delete(notificationID)
}

// GetUnreadNotificationCount возвращает количество непрочитанных уведомлений
func GetUnreadNotificationCount(db *gorm.DB, userID uint) (int64, error) {
	repo := repositories.NewNotificationRepository(db)
	return repo.CountUnread(userID)
}
