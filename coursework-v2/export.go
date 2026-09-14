package main

import (
	"encoding/csv"
	"fmt"
	"io"

	"szi-registry/models"
	"szi-registry/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

func exportRecords(window fyne.Window, records []models.SZIRecord) {
	if len(records) == 0 {
		dialog.ShowInformation("Экспорт", "Нет записей для сохранения", window)
		return
	}

	fileDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(fmt.Errorf("не удалось открыть файл для сохранения: %v", err), window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()
		if err := writeRecordsCSV(writer, records); err != nil {
			dialog.ShowError(fmt.Errorf("не удалось сохранить экспорт: %v", err), window)
			return
		}
		dialog.ShowInformation("Экспорт завершён", "Записи сохранены в выбранный CSV-файл", window)
	}, window)
	fileDialog.SetFileName("реестр_сзи.csv")
	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
	fileDialog.Show()
}

func writeRecordsCSV(output io.Writer, records []models.SZIRecord) error {
	writer := csv.NewWriter(output)
	headers := []string{
		"Наименование", "Тип СЗИ", "Класс СВТ", "Уровень доверия", "Класс АС",
		"Заявитель или разработчик", "Версия", "Номер сертификата", "Схема сертификации",
		"Дата выдачи", "Срок действия", "Место установки", "Ответственный", "Статус", "Примечания",
	}
	if err := writer.Write(headers); err != nil {
		return err
	}
	for _, record := range records {
		row := []string{
			record.Name, referenceType(record), referenceProtectionClass(record), record.TrustLevel,
			record.ACClass, referenceVendor(record), valueOrFallback(record.Version, record.SoftwareVersion),
			referenceCertificate(record), record.CertificationScheme, referenceIssueDate(record).Format("2006-01-02"),
			referenceExpiryDate(record).Format("2006-01-02"), referenceLocation(record), record.ResponsiblePerson,
			services.CalculateSziStatus(referenceExpiryDate(record), now()), record.Notes,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
