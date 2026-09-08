package main

import (
	"fmt"
	"szi-registry/config"
	"szi-registry/models"
	"szi-registry/repositories"
	"szi-registry/ui"

	"fyne.io/fyne/v2/app"
)

func main() {

	fmt.Println("1. Загружаем конфигурацию...")
	cfg := config.LoadConfig()
	fmt.Printf("2. Конфигурация загружена: %+v\n", cfg)

	fmt.Println("3. Подключаемся к базе данных...")
	dbManager, err := repositories.NewDBManagerWithConfig(cfg)
	if err != nil {
		fmt.Printf("Ошибка подключения к базе данных: %v\n", err)
		return
	}
	defer dbManager.Close()
	fmt.Println("4. Успешно подключились к базе данных")

	dbManager.DB.AutoMigrate(&models.User{}, &models.SZIRecord{}, &models.ReadNotification{}, &models.SZIAuditLog{})

	fmt.Println("5. Создаем приложение...")
	myApp := app.NewWithID("szi.registry.app")
	ui.ApplyFSTECTheme(myApp)
	fmt.Println("6. Приложение создано")

	fmt.Println("7. Открываем окно входа...")
	ui.ShowLoginWindow(myApp, dbManager.DB)
	fmt.Println("8. Окно входа открыто")

	fmt.Println("9. Запускаем приложение...")
	myApp.Run()
	fmt.Println("10. Приложение завершено")
}
