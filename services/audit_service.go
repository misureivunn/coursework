package services

import (
	"encoding/json"
	"gorm.io/gorm"
	"szi-registry/models"
	"szi-registry/repositories"
	"time"
)

// LogAction записывает действие пользователя в журнал аудита
func LogAction(db *gorm.DB, userID, recordID uint, action, entityType string, oldValue, newValue interface{}, description string) error {
	repo := repositories.NewAuditLogRepository(db)

	oldJSON, _ := json.Marshal(oldValue)
	newJSON, _ := json.Marshal(newValue)

	log := &models.AuditLog{
		UserID:      userID,
		RecordID:    recordID,
		Action:      action,
		EntityType:  entityType,
		OldValue:    string(oldJSON),
		NewValue:    string(newJSON),
		Description: description,
		CreatedAt:   time.Now(),
	}

	return repo.Create(log)
}

// GetUserAuditLogs возвращает последние N записей аудита для пользователя
func GetUserAuditLogs(db *gorm.DB, userID uint, limit int) ([]models.AuditLog, error) {
	repo := repositories.NewAuditLogRepository(db)
	return repo.FindByUserID(userID, limit)
}

// GetRecordAuditLogs возвращает историю изменений конкретной записи
func GetRecordAuditLogs(db *gorm.DB, recordID uint) ([]models.AuditLog, error) {
	repo := repositories.NewAuditLogRepository(db)
	return repo.FindByRecordID(recordID)
}

// GetAllAuditLogs возвращает все записи аудита (для админа)
func GetAllAuditLogs(db *gorm.DB, limit int) ([]models.AuditLog, error) {
	repo := repositories.NewAuditLogRepository(db)
	return repo.FindAll(limit)
}
