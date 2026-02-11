package main

import (
	"fmt"
	"szi-registry/config"
	"szi-registry/repositories"
	"szi-registry/services"
	"szi-registry/ui"
	"time"

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

	// Генерируем уведомления при запуске
	fmt.Println("5. Генерируем уведомления об истекающих сертификатах...")
	if err := services.GenerateExpiryNotifications(dbManager.DB); err != nil {
		fmt.Printf("Предупреждение: не удалось сгенерировать уведомления: %v\n", err)
	}

	// Периодически проверяем (каждые 24 часа)
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			fmt.Println("Периодическая генерация уведомлений...")
			if err := services.GenerateExpiryNotifications(dbManager.DB); err != nil {
				fmt.Printf("Предупреждение: не удалось сгенерировать уведомления: %v\n", err)
			}
		}
	}()

	fmt.Println("6. Создаем приложение...")
	myApp := app.NewWithID("szi.registry.app")
	fmt.Println("7. Приложение создано")

	fmt.Println("8. Открываем окно входа...")
	ui.ShowLoginWindow(myApp, dbManager.DB)
	fmt.Println("9. Окно входа открыто")

	fmt.Println("10. Запускаем приложение...")
	myApp.Run()
	fmt.Println("11. Приложение завершено")
}
