package utils

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func TestPasswordHash(hashedPassword, originalPassword string) {
	err := ValidatePassword(hashedPassword, originalPassword)
	if err != nil {
		log.Fatal("Пароли не совпадают:", err)
	}

	fmt.Println("Пароли совпадают! Вход должен работать.")
}
