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

// импортирует записи СЗИ из CSV файла
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

	// Поддерживаем два формата заголовков: со старыми (с ID) и новыми (без ID)
	var expectedHeaders []string
	var hasIDColumn bool

	if len(headers) == 15 { // со столбцом ID
		expectedHeaders = []string{
			"ID", "Наименование", "Тип", "Номер сертификата",
			"Дата выдачи", "Срок действия", "Место установки",
			"Статус", "Производитель", "Версия ПО", "Назначение",
			"Тип развертывания", "Класс защищенности", "Дата создания", "Дата обновления",
		}
		hasIDColumn = true
	} else if len(headers) == 14 { // без столбца ID
		expectedHeaders = []string{
			"Наименование", "Тип", "Номер сертификата",
			"Дата выдачи", "Срок действия", "Место установки",
			"Статус", "Производитель", "Версия ПО", "Назначение",
			"Тип развертывания", "Класс защищенности", "Дата создания", "Дата обновления",
		}
		hasIDColumn = false
	} else {
		return nil, fmt.Errorf("некорректное количество столбцов в CSV файле: ожидается 14 или 15 столбцов, получено %d", len(headers))
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

		// Парсим дату выдачи
		issueDateIndex := 3
		if hasIDColumn {
			issueDateIndex = 4
		}
		issueDate, err := time.Parse("2006-01-02", row[issueDateIndex])
		if err != nil {
			return nil, fmt.Errorf("ошибка при парсинге даты выдачи в строке %d: %v", rowIndex+2, err)
		}

		// Парсим дату истечения
		expiryDateIndex := 4
		if hasIDColumn {
			expiryDateIndex = 5
		}
		expiryDate, err := time.Parse("2006-01-02", row[expiryDateIndex])
		if err != nil {
			return nil, fmt.Errorf("ошибка при парсинге срока действия в строке %d: %v", rowIndex+2, err)
		}

		// Определяем статус на основе даты истечения
		statusIndex := 6
		if hasIDColumn {
			statusIndex = 7
		}
		status := row[statusIndex]
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
		var sziRecord models.SZIRecord

		if hasIDColumn {
			// Парсим ID (может быть пустым для новых записей)
			var id uint
			if row[0] != "" {
				parsedID, err := strconv.Atoi(row[0])
				if err != nil {
					return nil, fmt.Errorf("ошибка при парсинге ID в строке %d: %v", rowIndex+2, err)
				}
				id = uint(parsedID)
			}

			sziRecord = models.SZIRecord{
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

				// Поля для классификации СЗИ от НСД
				Purpose:         strings.TrimSpace(row[10]),
				DeploymentType:  strings.TrimSpace(row[11]),
				ClassProtection: strings.TrimSpace(row[12]),
			}
		} else {
			// Если нет ID в файле, создаем запись без него
			sziRecord = models.SZIRecord{
				Name:              strings.TrimSpace(row[0]),
				Type:              strings.TrimSpace(row[1]),
				CertNumber:        strings.TrimSpace(row[2]),
				CertIssueDate:     issueDate,
				CertExpiryDate:    expiryDate,
				Location:          strings.TrimSpace(row[5]),
				Status:            status,
				UserID:            userID, // Привязываем к текущему пользователю
				Manufacturer:      strings.TrimSpace(row[7]),
				SoftwareVersion:   strings.TrimSpace(row[8]),

				// Поля для классификации СЗИ от НСД
				Purpose:         strings.TrimSpace(row[9]),
				DeploymentType:  strings.TrimSpace(row[10]),
				ClassProtection: strings.TrimSpace(row[11]),
			}
		}

		// Парсим даты создания и обновления, если они указаны
		createdAtIndex := 12
		if hasIDColumn {
			createdAtIndex = 13
		}
		if len(row) > createdAtIndex && row[createdAtIndex] != "" {
			createdAt, err := parseDateTime(row[createdAtIndex])
			if err != nil {
				return nil, fmt.Errorf("ошибка при парсинге даты создания в строке %d: %v", rowIndex+2, err)
			}
			sziRecord.CreatedAt = createdAt
		} else {
			sziRecord.CreatedAt = time.Now()
		}

		updatedAtIndex := 13
		if hasIDColumn {
			updatedAtIndex = 14
		}
		if len(row) > updatedAtIndex && row[updatedAtIndex] != "" {
			updatedAt, err := parseDateTime(row[updatedAtIndex])
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

// пытается распознать дату в различных форматах
func parseDateTime(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("пустая строка даты")
	}

	// Основные форматы дат, которые могут использоваться
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"02.01.2006 15:04:05", 
		"02/01/2006 15:04:05", 
		"2006-01-02",          
		"02.01.2006",          
		"02/01/2006",          
		"02-Jan-2006",         
		"02-Jan-2006 15:04:05", 
		"Jan 02, 2006",       
		"Jan 02, 2006 15:04:05", 
		"1",             
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	// Если строка содержит только число, предполагаем, что это год
	if year, err := strconv.Atoi(dateStr); err == nil && year > 1900 && year < 2100 {
		return time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC), nil
	}

	return time.Time{}, fmt.Errorf("не удалось распознать формат даты: %s", dateStr)
}
