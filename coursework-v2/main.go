package main

import (
	"fmt"

	"coursework-v2/internal/application"
	"coursework-v2/internal/config"
	"coursework-v2/internal/storage"

	"fyne.io/fyne/v2/app"
)

func main() {
	cfg := config.Load()
	dsn, err := cfg.DSN()
	if err != nil {
		fmt.Printf("Ошибка подключения к базе данных: %v\n", err)
		fmt.Println("Запустите coursework-v2 PostgreSQL и задайте DB_PASSWORD и DB_PORT=5434 перед запуском, например:")
		fmt.Println("Задайте DB_PASSWORD и DB_PORT в локальном окружении и повторите запуск")
		return
	}
	db, err := storage.Open(dsn)
	if err != nil {
		fmt.Printf("Ошибка подключения к базе данных: %v\n", err)
		return
	}

	if err := storage.Migrate(db); err != nil {
		fmt.Printf("Ошибка подготовки базы данных: %v\n", err)
		return
	}
	if err := storage.NormalizeLegacyFields(db); err != nil {
		fmt.Printf("Не удалось обновить поля реестра: %v\n", err)
	}

	myApp := app.NewWithID("szi.registry.coursework.v2")
	applyCourseworkTheme(myApp)
	showLogin(myApp, application.New(db))
	myApp.Run()
}
