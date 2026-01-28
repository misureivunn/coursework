package services

import (
	"gorm.io/gorm"
	"szi-registry/csvimport"
	"szi-registry/repositories"
)

// ImportSziRecordsFromCSV импортирует записи СЗИ из CSV файла
func ImportSziRecordsFromCSV(db *gorm.DB, userID uint, filename string) (int, error) {
	// Импортируем записи из CSV
	records, err := csvimport.ImportFromCSV(filename, userID)
	if err != nil {
		return 0, err
	}

	// Сохраняем записи в базу данных
	repo := repositories.NewSZIRecordRepository(db)
	createdCount := 0

	for _, record := range records {
		// Проверяем, существует ли запись с таким ID
		_, err := repo.FindByID(record.ID)

		if err != nil && err == gorm.ErrRecordNotFound {
			// Если записи не существует, создаем новую
			err = repo.Create(&record)
			if err != nil {
				// Продолжаем импорт, даже если одна из записей не сохранилась
				continue
			}
			createdCount++
		} else if err == nil {
			// Если запись существует, обновляем её
			err = repo.Update(&record)
			if err != nil {
				// Продолжаем импорт, даже если одна из записей не обновилась
				continue
			}
			createdCount++
		}
	}

	return createdCount, nil
}
