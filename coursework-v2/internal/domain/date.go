package domain

import (
	"fmt"
	"time"
)

func ParseDate(value string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("некорректная дата «%s». Используйте формат ГГГГ-ММ-ДД", value)
	}
	return parsed, nil
}
