package application

import (
	"errors"
	"time"

	"coursework-v2/internal/domain"
	"coursework-v2/internal/storage"

	"gorm.io/gorm"
)

type App struct{ db *gorm.DB }

func New(db *gorm.DB) App { return App{db: db} }

func (a App) Login(username, password string) (domain.User, error) {
	return storage.Authenticate(a.db, username, password)
}
func (a App) Register(username, password string) error {
	return storage.CreateUser(a.db, username, password)
}
func (a App) Records(userID uint) ([]domain.SZIRecord, error) {
	return storage.ListRecords(a.db, userID)
}
func (a App) Status(expiry time.Time) string {
	return domain.StatusFor(expiry, time.Now())
}
func (a App) Create(record domain.SZIRecord) error {
	if err := storage.ValidateRecord(record); err != nil {
		return err
	}
	record.Status = storage.CalculateStatus(record.ExpiryDate, time.Now())
	return storage.CreateRecord(a.db, record)
}
func (a App) Update(record domain.SZIRecord) error {
	if err := storage.ValidateRecord(record); err != nil {
		return err
	}
	record.Status = storage.CalculateStatus(record.ExpiryDate, time.Now())
	return storage.UpdateRecord(a.db, record)
}
func (a App) Delete(userID, recordID uint) error {
	if err := storage.DeleteRecord(a.db, userID, recordID); err != nil {
		return errors.New("не удалось удалить запись")
	}
	return nil
}
