package services

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"szi-registry/models"
	"szi-registry/repositories"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// проверяет учетные данные пользователя
func AuthenticateUser(db *gorm.DB, username, password string) (*models.User, error) {
	repo := repositories.NewUserRepository(db)
	user, err := repo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gorm.ErrRecordNotFound
	}

	if !CheckPassword(password, user.Password) {
		return nil, gorm.ErrRecordNotFound
	}

	return user, nil
}

// создает нового пользователя
func CreateUser(db *gorm.DB, username, password string) (uint, error) {
	// Хешируем пароль
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return 0, err
	}

	user := &models.User{
		Username: username,
		Password: hashedPassword,
		Role:     "user",
	}

	repo := repositories.NewUserRepository(db)
	err = repo.Create(user)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
