package repositories

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"szi-registry/models"
)

type SZIRecordRepository struct {
	DB *gorm.DB
}

func NewSZIRecordRepository(db *gorm.DB) *SZIRecordRepository {
	return &SZIRecordRepository{DB: db}
}

func (r *SZIRecordRepository) Create(record *models.SZIRecord) error {
	return r.DB.Create(record).Error
}

func (r *SZIRecordRepository) FindByID(id uint) (*models.SZIRecord, error) {
	var record models.SZIRecord
	err := r.DB.Preload("User").First(&record, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &record, err
}

func (r *SZIRecordRepository) FindByUserID(userID uint) ([]models.SZIRecord, error) {
	var records []models.SZIRecord
	err := r.DB.Where("user_id = ?", userID).Find(&records).Error
	return records, err
}

func (r *SZIRecordRepository) FindAccessibleByUserID(userID uint) ([]models.SZIRecord, error) {
	var records []models.SZIRecord

	err := r.DB.Raw(`
		SELECT DISTINCT sr.*
		FROM szi_records sr
		LEFT JOIN access_permissions ap ON sr.id = ap.record_id
		WHERE sr.user_id = ? OR ap.guest_user_id = ?
		ORDER BY sr.name
	`, userID, userID).Scan(&records).Error

	return records, err
}

func (r *SZIRecordRepository) Update(record *models.SZIRecord) error {
	return r.DB.Save(record).Error
}

func (r *SZIRecordRepository) Delete(id uint) error {
	return r.DB.Delete(&models.SZIRecord{}, id).Error
}

func (r *SZIRecordRepository) Search(userID uint, criteria map[string]interface{}) ([]models.SZIRecord, error) {
	var records []models.SZIRecord

	query := r.DB.Model(&models.SZIRecord{}).
		Joins("LEFT JOIN access_permissions ap ON szi_records.id = ap.record_id").
		Where("szi_records.user_id = ? OR ap.guest_user_id = ?", userID, userID)

	for field, value := range criteria {
		switch field {
		case "name":
			query = query.Where("szi_records.name ILIKE ?", "%"+value.(string)+"%")
		case "type":
			query = query.Where("szi_records.type = ?", value)
		case "location":
			query = query.Where("szi_records.location ILIKE ?", "%"+value.(string)+"%")
		case "status":
			query = query.Where("szi_records.status = ?", value)
		case "cert_number":
			query = query.Where("szi_records.cert_number ILIKE ?", "%"+value.(string)+"%")
		case "manufacturer":
			query = query.Where("szi_records.manufacturer ILIKE ?", "%"+value.(string)+"%")
		}
	}

	err := query.Order("szi_records.name").Find(&records).Error
	return records, err
}

func (r *SZIRecordRepository) DeleteForUser(userID, recordID uint) error {
	// Проверяем, принадлежит ли запись пользователю
	var record models.SZIRecord
	err := r.DB.Where("id = ? AND user_id = ?", recordID, userID).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("запись не найдена или не принадлежит пользователю")
		}
		return err
	}

	// Удаляем запись
	return r.DB.Delete(&record).Error
}
