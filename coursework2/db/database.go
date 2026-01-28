package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase(databasePath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}