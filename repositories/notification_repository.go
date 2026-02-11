package repositories

import (
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

func (r *NotificationRepository) FindByUserID(userID uint, unreadOnly bool) ([]models.Notification, error) {
	var notifications []models.Notification
	query := r.DB.Where("user_id = ?", userID)
	
	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}
	
	err := query.Order("priority DESC, created_at DESC").
		Preload("Record").
		Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) MarkAsRead(id uint) error {
	return r.DB.Model(&models.Notification{}).
		Where("id = ?", id).
		Update("is_read", true).Error
}

func (r *NotificationRepository) MarkAllAsRead(userID uint) error {
	return r.DB.Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error
}

func (r *NotificationRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Notification{}, id).Error
}

func (r *NotificationRepository) CountUnread(userID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}
