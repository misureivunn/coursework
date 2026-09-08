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
		return
	}
	defer dbManager.Close()

	dbManager.DB.AutoMigrate(&models.User{}, &models.SZIRecord{}, &models.ReadNotification{}, &models.SZIAuditLog{})

	myApp := app.NewWithID("szi.registry.coursework.v2")
	applyCourseworkTheme(myApp)
	showLogin(myApp, dbManager.DB)
	myApp.Run()
}
