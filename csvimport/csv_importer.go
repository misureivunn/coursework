package csvimport

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"szi-registry/models"
	"time"
)

// ImportFromCSV импортирует записи СЗИ из CSV файла
func ImportFromCSV(filename string, userID uint) ([]models.SZIRecord, error) {
	// Открываем файл
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии файла: %v", err)
	}
	defer file.Close()

	// Создаем CSV ридер
	reader := csv.NewReader(file)

	// Читаем все строки
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении CSV файла: %v", err)
	}

	// Проверяем, есть ли заголовки
	if len(records) < 2 {
		return nil, fmt.Errorf("файл CSV должен содержать хотя бы одну строку данных")
	}

	// Проверяем заголовки
	headers := records[0]
	expectedHeaders := []string{
		"ID", "Наименование", "Тип", "Номер сертификата",
		"Дата выдачи", "Срок действия", "Место установки",
		"Статус", "Производитель", "Версия ПО", "Контактное лицо",
		"Ссылка на документацию", "Дата создания", "Дата обновления",
	}

	if len(headers) != len(expectedHeaders) {
		return nil, fmt.Errorf("некорректное количество столбцов в CSV файле")
	}

	// Сравниваем заголовки (без учета регистра и пробелов)
	for i, header := range headers {
		if strings.TrimSpace(strings.ToLower(header)) != strings.ToLower(expectedHeaders[i]) {
			return nil, fmt.Errorf("ожидаемый заголовок '%s', получен '%s'", expectedHeaders[i], header)
		}
	}

	var importedRecords []models.SZIRecord

	// Обрабатываем строки данных (начиная с индекса 1, т.к. 0 - заголовки)
	for rowIndex, row := range records[1:] {
		if len(row) != len(expectedHeaders) {
			return nil, fmt.Errorf("некорректное количество столбцов в строке %d", rowIndex+2)
		}

		// Парсим ID (может быть пустым для новых записей)
		var id uint
		if row[0] != "" {
			parsedID, err := strconv.Atoi(row[0])
			if err != nil {
				return nil, fmt.Errorf("ошибка при парсинге ID в строке %d: %v", rowIndex+2, err)
			}
			id = uint(parsedID)
		}

		// Парсим дату выдачи
		issueDate, err := time.Parse("2006-01-02", row[4])
		if err != nil {
			return nil, fmt.Errorf("ошибка при парсинге даты выдачи в строке %d: %v", rowIndex+2, err)
		}

		// Парсим дату истечения
		expiryDate, err := time.Parse("2006-01-02", row[5])
		if err != nil {
			return nil, fmt.Errorf("ошибка при парсинге срока действия в строке %d: %v", rowIndex+2, err)
		}

		// Определяем статус на основе даты истечения
		status := row[7]
		if status == "" {
			now := time.Now()
			if expiryDate.Before(now) {
				status = "Просрочено"
			} else {
				// Проверяем, если срок истекает в ближайшие 30 дней
				in30Days := now.AddDate(0, 0, 30)
				if expiryDate.Before(in30Days) {
					status = "Скоро истекает"
				} else {
					status = "Актуально"
				}
			}
		}

		// Создаем запись СЗИ
		sziRecord := models.SZIRecord{
			ID:                id,
			Name:              strings.TrimSpace(row[1]),
			Type:              strings.TrimSpace(row[2]),
			CertNumber:        strings.TrimSpace(row[3]),
			CertIssueDate:     issueDate,
			CertExpiryDate:    expiryDate,
			Location:          strings.TrimSpace(row[6]),
			Status:            status,
			UserID:            userID, // Привязываем к текущему пользователю
			Manufacturer:      strings.TrimSpace(row[8]),
			SoftwareVersion:   strings.TrimSpace(row[9]),
			ContactPerson:     strings.TrimSpace(row[10]),
			DocumentationLink: strings.TrimSpace(row[11]),
		}

		// Парсим даты создания и обновления, если они указаны
		if row[12] != "" {
			createdAt, err := time.Parse("2006-01-02 15:04:05", row[12])
			if err != nil {
				return nil, fmt.Errorf("ошибка при парсинге даты создания в строке %d: %v", rowIndex+2, err)
			}
			sziRecord.CreatedAt = createdAt
		} else {
			sziRecord.CreatedAt = time.Now()
		}

		if row[13] != "" {
			updatedAt, err := time.Parse("2006-01-02 15:04:05", row[13])
			if err != nil {
				return nil, fmt.Errorf("ошибка при парсинге даты обновления в строке %d: %v", rowIndex+2, err)
			}
			sziRecord.UpdatedAt = updatedAt
		} else {
			sziRecord.UpdatedAt = time.Now()
		}

		importedRecords = append(importedRecords, sziRecord)
	}

	return importedRecords, nil
}
