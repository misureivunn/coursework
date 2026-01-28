package main

import (
	"fmt"
	"szi-registry-new/config"
	"szi-registry-new/db"
	"szi-registry-new/models"
	"szi-registry-new/ui"

	"fyne.io/fyne/v2/app"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Запуск реестра СЗИ...")

	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Инициализируем базу данных
	database, err := db.InitDatabase(cfg.DatabasePath)
	if err != nil {
		fmt.Printf("Ошибка инициализации базы данных: %v\n", err)
		return
	}

	// Миграция моделей
	err = database.AutoMigrate(&models.User{}, &models.SZI{})
	if err != nil {
		fmt.Printf("Ошибка миграции базы данных: %v\n", err)
		return
	}

	// Проверяем наличие администратора, если нет - создаем
	createAdminUser(database)

	// Запускаем GUI приложение
	a := app.New()
	ui.ShowLoginScreen(a, database)
}

// createAdminUser создает пользователя admin при первом запуске, если он не существует
func createAdminUser(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)
	
	if count == 0 {
		password := "admin123" // Пароль по умолчанию для первого запуска
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("Ошибка хеширования пароля: %v\n", err)
			return
		}

		adminUser := models.User{
			Username: "admin",
			Password: string(hashedPassword),
			Role:     "admin",
		}

		result := db.Create(&adminUser)
		if result.Error != nil {
			fmt.Printf("Ошибка создания администратора: %v\n", result.Error)
		} else {
			fmt.Println("Создан пользователь admin с паролем admin123")
		}
	}
}