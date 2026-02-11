package repositories

import (
	"gorm.io/gorm"
	"szi-registry/models"
)

type AuditLogRepository struct {
	DB *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{DB: db}
}

func (r *AuditLogRepository) Create(log *models.AuditLog) error {
	return r.DB.Create(log).Error
}

func (r *AuditLogRepository) FindByUserID(userID uint, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Preload("User").
		Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) FindByRecordID(recordID uint) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.DB.Where("record_id = ?", recordID).
		Order("created_at DESC").
		Preload("User").
		Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) FindAll(limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.DB.Order("created_at DESC").
		Limit(limit).
		Preload("User").
		Find(&logs).Error
	return logs, err
}
