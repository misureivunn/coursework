package main

import (
	"fmt"

	"szi-registry/config"
	"szi-registry/models"
	"szi-registry/repositories"

	"fyne.io/fyne/v2/app"
)

func main() {
	cfg := config.LoadConfig()
	dbManager, err := repositories.NewDBManagerWithConfig(cfg)
	if err != nil {
		fmt.Printf("Ошибка подключения к базе данных: %v\n", err)
		fmt.Println("Запустите coursework-v2 PostgreSQL и задайте DB_PASSWORD и DB_PORT=5434 перед запуском, например:")
		fmt.Println("DB_PASSWORD=local-dev-password DB_PORT=5434 go run ./coursework-v2")
		return
	}
	defer dbManager.Close()

	dbManager.DB.AutoMigrate(&models.User{}, &models.SZIRecord{}, &models.ReadNotification{}, &models.SZIAuditLog{})

	myApp := app.NewWithID("szi.registry.coursework.v2")
	applyCourseworkTheme(myApp)
	showLogin(myApp, dbManager.DB)
	myApp.Run()
}
