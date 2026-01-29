package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
	"szi-registry/models"
)

// ShowErrorDialog отображает диалог с сообщением об ошибке
func ShowErrorDialog(message string, canvas fyne.Canvas) {
	var dialog *widget.PopUp

	content := widget.NewRichTextFromMarkdown("# Ошибка\n\n" + message)
	content.Wrapping = fyne.TextWrapWord

	dialog = widget.NewModalPopUp(
		container.NewVBox(
			content,
			widget.NewButton("Закрыть", func() {
				dialog.Hide()
			}),
		),
		canvas,
	)
	dialog.Show()
}

// ShowSuccessDialog отображает диалог с сообщением об успешном выполнении
func ShowSuccessDialog(message string, canvas fyne.Canvas) {
	var dialog *widget.PopUp

	content := widget.NewRichTextFromMarkdown("# Успешно\n\n" + message)
	content.Wrapping = fyne.TextWrapWord

	dialog = widget.NewModalPopUp(
		container.NewVBox(
			content,
			widget.NewButton("Закрыть", func() {
				dialog.Hide()
			}),
		),
		canvas,
	)
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
	status := record.Status
	expiryDate := record.CertExpiryDate
	now := time.Now()

	// Если статус не установлен явно, определяем его по дате истечения
	if status == "" {
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
	return status
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