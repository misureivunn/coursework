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


func (r *SZIRecordRepository) Update(record *models.SZIRecord) error {
	return r.DB.Save(record).Error
}

func (r *SZIRecordRepository) Delete(id uint) error {
	return r.DB.Delete(&models.SZIRecord{}, id).Error
}

func (r *SZIRecordRepository) Search(userID uint, criteria map[string]interface{}) ([]models.SZIRecord, error) {
	var records []models.SZIRecord

	// Основной запрос с фильтрацией
	query := r.DB.Where("user_id = ?", userID)

	// Применяем дополнительные критерии поиска
	for field, value := range criteria {
		switch field {
		case "name":
			query = query.Where("name ILIKE ?", "%"+value.(string)+"%")
		case "type":
			query = query.Where("szi_type = ?", value)
		case "location":
			query = query.Where("location ILIKE ?", "%"+value.(string)+"%")
		case "status":
			query = query.Where("status = ?", value)
		case "cert_number":
			query = query.Where("cert_number ILIKE ?", "%"+value.(string)+"%")
		case "manufacturer":
			query = query.Where("manufacturer ILIKE ?", "%"+value.(string)+"%")
		case "purpose":
			query = query.Where("purpose = ?", value)
		case "deployment_type":
			query = query.Where("deployment_type = ?", value)
		case "class_protection":
			query = query.Where("class_protection = ?", value)
		}
	}

	err := query.Order("name").Find(&records).Error
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