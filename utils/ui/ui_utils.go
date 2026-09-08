package ui

import (
	"szi-registry/models"
	"szi-registry/services"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowErrorDialog(message string, canvas fyne.Canvas) {
	var dialog *widget.PopUp

	content := widget.NewRichTextFromMarkdown("# Ошибка\n\n" + message)
	content.Wrapping = fyne.TextWrapWord

	// Создаём широкий контейнер
	contentContainer := container.NewVBox(
		content,
		widget.NewButton("Закрыть", func() {
			dialog.Hide()
		}),
	)

	dialog = widget.NewModalPopUp(contentContainer, canvas)

	// Устанавливаем размер: широкий (600px) и низкий (150px)
	dialog.Resize(fyne.NewSize(600, 150))

	dialog.Show()
}

func ShowSuccessDialog(message string, canvas fyne.Canvas) {
	var dialog *widget.PopUp

	content := widget.NewRichTextFromMarkdown("# Успешно\n\n" + message)
	content.Wrapping = fyne.TextWrapWord

	contentContainer := container.NewVBox(
		content,
		widget.NewButton("Закрыть", func() {
			dialog.Hide()
		}),
	)

	dialog = widget.NewModalPopUp(contentContainer, canvas)

	// Устанавливаем размер: широкий и низкий
	dialog.Resize(fyne.NewSize(600, 150))

	dialog.Show()
}

// ParseDate parses date string in format YYYY-MM-DD
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// FormatDate formats date as YYYY-MM-DD
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// GetRecordStatus определяет статус записи на основе даты истечения
func GetRecordStatus(record models.SZIRecord) string {
	return services.CalculateSziStatus(record.CertExpiryDate, time.Now())
}

// GetCertificateStats возвращает статистику по сертификатам
func GetCertificateStats(records []models.SZIRecord) (active, expired, noCert int) {
	now := time.Now()

	for _, record := range records {
		if record.CertExpiryDate.IsZero() {
			noCert++
		} else if record.CertExpiryDate.Before(now) {
			expired++
		} else {
			active++
		}
	}

	return active, expired, noCert
}
