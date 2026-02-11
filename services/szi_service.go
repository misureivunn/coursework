package services

import (
	"fmt"
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
	err := repo.Create(record)

	if err == nil {
		// Логируем создание
		LogAction(db, record.UserID, record.ID, "CREATE", "SZIRecord", nil, record, fmt.Sprintf("Создана запись СЗИ: %s", record.Name))
	}

	return err
}

func UpdateSziRecord(db *gorm.DB, record *models.SZIRecord) error {
	repo := repositories.NewSZIRecordRepository(db)

	// Получаем старое значение
	oldRecord, _ := repo.FindByID(record.ID)

	err := repo.Update(record)

	if err == nil {
		// Логируем изменение
		LogAction(db, record.UserID, record.ID, "UPDATE", "SZIRecord", oldRecord, record, fmt.Sprintf("Обновлена запись СЗИ: %s", record.Name))
	}

	return err
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

	// Получаем запись перед удалением
	record, _ := repo.FindByID(recordID)

	// Проверяем, принадлежит ли запись пользователю
	var existingRecord models.SZIRecord
	err := db.Where("id = ? AND user_id = ?", recordID, userID).First(&existingRecord).Error
	if err != nil {
		return err
	}

	err = repo.Delete(recordID)

	if err == nil && record != nil {
		// Логируем удаление
		LogAction(db, userID, recordID, "DELETE", "SZIRecord", record, nil, fmt.Sprintf("Удалена запись СЗИ: %s", record.Name))
	}

	return err
}

// возвращает всех пользователей
func GetAllUsers(db *gorm.DB) ([]models.User, error) {
	repo := repositories.NewUserRepository(db)
	return repo.GetAll()
}
