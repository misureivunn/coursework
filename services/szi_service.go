package services

import (
	"gorm.io/gorm"
	"szi-registry/models"
	"szi-registry/repositories"
)

func GetSziRecordsByUserID(db *gorm.DB, userID uint) ([]models.SZIRecord, error) {
	repo := repositories.NewSZIRecordRepository(db)
	return repo.FindByUserID(userID)
}

func AddSziRecord(db *gorm.DB, record *models.SZIRecord) error {
	repo := repositories.NewSZIRecordRepository(db)
	return repo.Create(record)
}

func UpdateSziRecord(db *gorm.DB, record *models.SZIRecord) error {
	repo := repositories.NewSZIRecordRepository(db)
	return repo.Update(record)
}

func GetSziRecordByID(db *gorm.DB, id uint) (*models.SZIRecord, error) {
	repo := repositories.NewSZIRecordRepository(db)
	return repo.FindByID(id)
}

// выполняет поиск записей СЗИ по различным критериям
func SearchSziRecords(db *gorm.DB, userID uint, filters map[string]interface{}) ([]models.SZIRecord, error) {
	repo := repositories.NewSZIRecordRepository(db)
	return repo.Search(userID, filters)
}

// возвращает все записи СЗИ, доступные пользователю
func GetSziRecordsAccessibleByUserID(db *gorm.DB, userID uint) ([]models.SZIRecord, error) {
	repo := repositories.NewSZIRecordRepository(db)
	return repo.FindAccessibleByUserID(userID)
}
