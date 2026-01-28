package services

import (
	"gorm.io/gorm"
	"szi-registry/export"
)

// экспортирует записи СЗИ в CSV файл
func ExportSziRecordsToCSV(db *gorm.DB, userID uint, filename string) error {
	// Получаем записи пользователя
	records, err := GetSziRecordsByUserID(db, userID)
	if err != nil {
		return err
	}

	// Экспортируем в CSV
	err = export.ExportToCSV(records, filename)
	if err != nil {
		return err
	}

	return nil
}
