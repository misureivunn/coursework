package services

import (
	"szi-registry/models"
	"time"

	"gorm.io/gorm"
)

// LogAudit создает запись в журнале аудита
func LogAudit(db *gorm.DB, sziID, userID uint, action, fieldName, oldValue, newValue, ipAddress string) error {
	auditLog := models.SZIAuditLog{
		SZIID:       sziID,
		UserID:      userID,
		Action:      action,
		FieldName:   fieldName,
		OldValue:    oldValue,
		NewValue:    newValue,
		PerformedAt: time.Now(),
		IPAddress:   ipAddress,
	}
	return db.Create(&auditLog).Error
}

// GetAuditLogsBySZI получает историю изменений для конкретной записи СЗИ
func GetAuditLogsBySZI(db *gorm.DB, sziID uint) ([]models.SZIAuditLog, error) {
	var logs []models.SZIAuditLog
	err := db.Where("szi_id = ?", sziID).
		Order("performed_at DESC").
		Preload("User").
		Find(&logs).Error
	return logs, err
}

// GetAuditLogsByUser получает историю действий конкретного пользователя
func GetAuditLogsByUser(db *gorm.DB, userID uint) ([]models.SZIAuditLog, error) {
	var logs []models.SZIAuditLog
	err := db.Where("user_id = ?", userID).
		Order("performed_at DESC").
		Preload("SZI").
		Find(&logs).Error
	return logs, err
}

// GetRecentAuditLogs получает последние записи аудита
func GetRecentAuditLogs(db *gorm.DB, limit int) ([]models.SZIAuditLog, error) {
	var logs []models.SZIAuditLog
	err := db.Order("performed_at DESC").
		Limit(limit).
		Preload("User").
		Preload("SZI").
		Find(&logs).Error
	return logs, err
}
