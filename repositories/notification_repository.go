package repositories

import (
	"errors"
	"gorm.io/gorm"
	"szi-registry/models"
)

type NotificationRepository struct {
	DB *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

func (r *NotificationRepository) Create(notification *models.Notification) error {
	return r.DB.Create(notification).Error
}

func (r *NotificationRepository) FindByRecipientID(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.DB.Where("recipient_id = ?", userID).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) MarkAsRead(notificationID uint) error {
	notification := models.Notification{}
	result := r.DB.First(&notification, notificationID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("уведомление не найдено")
		}
		return result.Error
	}

	notification.ReadStatus = true
	notification.ReadAt = &notification.UpdatedAt
	return r.DB.Save(&notification).Error
}
