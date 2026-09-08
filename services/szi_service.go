package services

import (
	"fmt"
	"strings"
	"szi-registry/models"
	"szi-registry/repositories"
	"time"

	"gorm.io/gorm"
)

const (
	StatusCurrent   = "Актуально"
	StatusAttention = "Требует внимания"
	StatusExpired   = "Просрочено"
)

func CalculateSziStatus(expiryDate time.Time, now time.Time) string {
	if expiryDate.Before(now) {
		return StatusExpired
	}
	if expiryDate.Before(now.AddDate(0, 0, 60)) {
		return StatusAttention
	}
	return StatusCurrent
}

func ValidateSziRecord(record *models.SZIRecord, now time.Time) error {
	if strings.TrimSpace(record.Name) == "" {
		return fmt.Errorf("наименование СЗИ обязательно")
	}
	if strings.TrimSpace(record.CertNumber) == "" {
		return fmt.Errorf("номер сертификата обязателен")
	}
	if record.CertExpiryDate.IsZero() {
		return fmt.Errorf("срок действия сертификата обязателен")
	}
	if record.CertIssueDate.IsZero() {
		return fmt.Errorf("дата выдачи сертификата обязательна")
	}
	if record.CertIssueDate.After(now) {
		return fmt.Errorf("дата выдачи не может быть позже текущей даты")
	}
	if record.CertExpiryDate.Before(record.CertIssueDate) {
		return fmt.Errorf("срок действия не может быть раньше даты выдачи")
	}
	record.Status = CalculateSziStatus(record.CertExpiryDate, now)
	return nil
}

func GetSziRecordsByUserID(db *gorm.DB, userID uint) ([]models.SZIRecord, error) {
	repo := repositories.NewSZIRecordRepository(db)
	return repo.FindByUserID(userID)
}

// AddSziRecord добавляет новую запись СЗИ
func AddSziRecord(db *gorm.DB, record *models.SZIRecord) error {
	return CreateSziRecord(db, record)
}

// CreateSziRecord создает новую запись СЗИ с логированием
func CreateSziRecord(db *gorm.DB, record *models.SZIRecord) error {
	if err := ValidateSziRecord(record, time.Now()); err != nil {
		return err
	}
	err := db.Create(record).Error
	if err != nil {
		return err
	}

	// Логируем создание записи
	_ = LogAudit(db, record.ID, record.UserID, "create", "", "", fmt.Sprintf("Создана запись: %s", record.Name), "")

	return nil
}

// UpdateSziRecord обновляет запись СЗИ с логированием изменений
func UpdateSziRecord(db *gorm.DB, currentUserID uint, record *models.SZIRecord) error {
	// Получаем старую версию записи для сравнения
	var oldRecord models.SZIRecord
	if err := db.First(&oldRecord, record.ID).Error; err != nil {
		return err
	}
	if oldRecord.UserID != currentUserID {
		return fmt.Errorf("запись не принадлежит пользователю")
	}
	record.UserID = currentUserID
	if err := ValidateSziRecord(record, time.Now()); err != nil {
		return err
	}

	// Обновляем запись
	err := db.Save(record).Error
	if err != nil {
		return err
	}

	// Логируем изменения
	if oldRecord.Name != record.Name {
		_ = LogAudit(db, record.ID, record.UserID, "update", "Name", oldRecord.Name, record.Name, "")
	}
	if oldRecord.Type != record.Type {
		_ = LogAudit(db, record.ID, record.UserID, "update", "Type", oldRecord.Type, record.Type, "")
	}
	if oldRecord.CertNumber != record.CertNumber {
		_ = LogAudit(db, record.ID, record.UserID, "update", "CertNumber", oldRecord.CertNumber, record.CertNumber, "")
	}
	if oldRecord.Status != record.Status {
		_ = LogAudit(db, record.ID, record.UserID, "update", "Status", oldRecord.Status, record.Status, "")
	}

	return nil
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

// DeleteSziRecord удаляет запись СЗИ с логированием
func DeleteSziRecord(db *gorm.DB, userID, recordID uint) error {
	var record models.SZIRecord
	if err := db.First(&record, recordID).Error; err != nil {
		return err
	}

	// Проверяем, принадлежит ли запись пользователю
	if record.UserID != userID {
		return fmt.Errorf("запись не принадлежит пользователю")
	}

	// Логируем удаление
	_ = LogAudit(db, record.ID, userID, "delete", "", fmt.Sprintf("Удалена запись: %s", record.Name), "", "")

	// Удаляем запись
	return db.Delete(&record).Error
}

// возвращает всех пользователей
func GetAllUsers(db *gorm.DB) ([]models.User, error) {
	repo := repositories.NewUserRepository(db)
	return repo.GetAll()
}
