package utils

import (
	"fmt"
	"log"
)

// логирует ошибку с дополнительной информацией
func LogError(operation string, err error) {
	if err != nil {
		log.Printf("ERROR in %s: %v", operation, err)
	}
}

// обрабатывает ошибку и возвращает форматированное сообщение
func HandleError(operation string, err error) error {
	if err != nil {
		log.Printf("ERROR in %s: %v", operation, err)
		return fmt.Errorf("ошибка в операции %s: %v", operation, err)
	}
	return nil
}

// проверяет данные перед экспортом
func ValidateExportData(data interface{}) bool {
	if data == nil {
		return false
	}
	return true
}
