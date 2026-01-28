package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"szi-registry/models"
	"szi-registry/utils"
)

// экспортирует записи СЗИ в CSV файл
func ExportToCSV(records []models.SZIRecord, filename string) error {
	// Проверяем данные перед экспортом
	if !utils.ValidateSZIRecordSlice(records) {
		return utils.HandleError("ValidateSZIRecordSlice", fmt.Errorf("данные для экспорта пусты"))
	}

	file, err := os.Create(filename)
	if err != nil {
		return utils.HandleError("os.Create", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Записываем заголовки
	headers := []string{
		"ID", "Наименование", "Тип", "Номер сертификата",
		"Дата выдачи", "Срок действия", "Место установки",
		"Статус", "Производитель", "Версия ПО", "Контактное лицо",
		"Ссылка на документацию", "Дата создания", "Дата обновления",
	}

	err = writer.Write(headers)
	if err != nil {
		return utils.HandleError("writer.Write headers", err)
	}

	// Записываем данные
	for _, record := range records {
		row := []string{
			strconv.Itoa(int(record.ID)), // ID
			record.Name,                  // Наименование
			record.Type,                  // Тип
			record.CertNumber,            // Номер сертификата
			record.CertIssueDate.Format("2006-01-02"),  // Дата выдачи
			record.CertExpiryDate.Format("2006-01-02"), // Срок действия
			record.Location,                                // Место установки
			record.Status,                                  // Статус
			record.Manufacturer,                            // Производитель
			record.SoftwareVersion,                         // Версия ПО
			record.ContactPerson,                           // Контактное лицо
			record.DocumentationLink,                       // Ссылка на документацию
			record.CreatedAt.Format("2006-01-02 15:04:05"), // Дата создания
			record.UpdatedAt.Format("2006-01-02 15:04:05"), // Дата обновления
		}

		err = writer.Write(row)
		if err != nil {
			return utils.HandleError("writer.Write row", err)
		}
	}

	return nil
}
