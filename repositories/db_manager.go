package repositories

import (
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
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5433",
		DBUser:     "szi_user",
		DBPassword: "MyV3ryS3cur3P@ss2026!",
		DBName:     "szi_registry",
	}
	return NewDBManagerWithConfig(cfg)
}

// создает новый менеджер базы данных с указанной конфигурацией
func NewDBManagerWithConfig(cfg *config.Config) (*DBManager, error) {
	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Автоматическая миграция схемы
	err = db.AutoMigrate(
		&models.User{},
		&models.SZIRecord{},
		&models.UserRole{},
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
