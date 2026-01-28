package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func PrintPasswordHash(password string) {
	hashedPassword, err := GeneratePasswordHash(password)
	if err != nil {
		fmt.Println("Ошибка при хешировании:", err)
		return
	}

	fmt.Println("Хеш для пароля '"+password+"':", hashedPassword)
}
