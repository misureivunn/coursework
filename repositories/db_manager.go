package repositories

import (
	"fmt"
	"szi-registry/config"
	"szi-registry/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// управляет соединениями с базой данных
type DBManager struct {
	DB *gorm.DB
}

// создает новый менеджер базы данных с настройками по умолчанию
func NewDBManager() (*DBManager, error) {
	return NewDBManagerWithConfig(config.LoadConfig())
}

// создает новый менеджер базы данных с указанной конфигурацией
func NewDBManagerWithConfig(cfg *config.Config) (*DBManager, error) {
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("переменная окружения DB_PASSWORD не задана")
	}
	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Автоматическая миграция схемы
	err = db.AutoMigrate(
		&models.User{},
		&models.SZIRecord{},
	)
	if err != nil {
		return nil, err
	}

	return &DBManager{DB: db}, nil
}

// закрывает соединение с базой данных
func (dm *DBManager) Close() error {
	sqlDB, err := dm.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
