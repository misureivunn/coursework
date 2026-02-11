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

// выполняет поиск записей СЗИ по различным критериям (альтернативное имя)
func SearchSZIRecords(db *gorm.DB, userID uint, filters map[string]interface{}) ([]models.SZIRecord, error) {
	repo := repositories.NewSZIRecordRepository(db)
	return repo.Search(userID, filters)
}

// возвращает пользователя по имени
func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	repo := repositories.NewUserRepository(db)
	return repo.FindByUsername(username)
}

// возвращает все записи СЗИ пользователя
func GetUserRecords(db *gorm.DB, userID uint) ([]models.SZIRecord, error) {
	repo := repositories.NewSZIRecordRepository(db)
	return repo.FindByUserID(userID)
}

// удаляет запись СЗИ
func DeleteSziRecord(db *gorm.DB, userID, recordID uint) error {
	repo := repositories.NewSZIRecordRepository(db)

	// Проверяем, принадлежит ли запись пользователю
	var record models.SZIRecord
	err := db.Where("id = ? AND user_id = ?", recordID, userID).First(&record).Error
	if err != nil {
		return err
	}

	return repo.Delete(recordID)
}

// возвращает всех пользователей
func GetAllUsers(db *gorm.DB) ([]models.User, error) {
	repo := repositories.NewUserRepository(db)
	return repo.GetAll()
}

