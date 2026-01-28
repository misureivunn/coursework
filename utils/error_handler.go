package utils

import (
	"fmt"
	"reflect"
	"szi-registry/models"
)

// HandleError обрабатывает ошибки и возвращает их
func HandleError(operation string, err error) error {
	if err != nil {
		return fmt.Errorf("ошибка в операции %s: %v", operation, err)
	}
	return nil
}

// ValidateExportData проверяет, есть ли данные для экспорта
func ValidateExportData(data interface{}) bool {
	if data == nil {
		return false
	}

	// Используем reflection для проверки, является ли это непустым срезом
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Slice {
		return val.Len() > 0
	}

	// Проверяем строку
	if str, ok := data.(string); ok {
		return len(str) > 0
	}

	return false
}

// ValidateSZIRecordSlice проверяет, есть ли данные для экспорта SZIRecord
func ValidateSZIRecordSlice(records []models.SZIRecord) bool {
	return len(records) > 0
}
