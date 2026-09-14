package storage

import (
	"errors"
	"strings"
	"time"

	"szi-registry/coursework-v2/internal/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type userRow struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"size:64;unique;not null"`
	Password string `gorm:"size:255;not null"`
	Role     string `gorm:"size:32;default:user"`
}

type recordRow struct {
	ID                  uint   `gorm:"primaryKey"`
	Name                string `gorm:"size:255;not null"`
	SZIType             string `gorm:"column:szi_type_reference;size:100"`
	ProtectionClass     string `gorm:"size:32"`
	TrustLevel          string `gorm:"size:10"`
	ACClass             string `gorm:"size:32"`
	Vendor              string `gorm:"size:255"`
	Version             string `gorm:"size:100"`
	CertificateNumber   string `gorm:"size:100"`
	CertificationScheme string `gorm:"size:50"`
	IssueDate           time.Time
	ExpiryDate          time.Time
	InstallLocation     string `gorm:"size:255"`
	ResponsiblePerson   string `gorm:"size:255"`
	Notes               string `gorm:"type:text"`
	Status              string `gorm:"size:20"`
	OwnerID             uint
	UserID              uint
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (userRow) TableName() string   { return "users" }
func (recordRow) TableName() string { return "szi_records" }

func Open(dsn string) (*gorm.DB, error) { return gorm.Open(postgres.Open(dsn), &gorm.Config{}) }
func Migrate(db *gorm.DB) error         { return db.AutoMigrate(&userRow{}, &recordRow{}) }

func NormalizeLegacyFields(db *gorm.DB) error {
	var rows []recordRow
	if err := db.Where("szi_type_reference = '' OR szi_type_reference IS NULL").Find(&rows).Error; err != nil {
		return err
	}
	for _, row := range rows {
		if row.SZIType == "" {
			row.SZIType = "Другое"
		}
		row.OwnerID = row.UserID
		if err := db.Save(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

func Authenticate(db *gorm.DB, username, password string) (domain.User, error) {
	var row userRow
	if err := db.Where("username = ?", username).First(&row).Error; err != nil {
		return domain.User{}, errors.New("неверное имя пользователя или пароль")
	}
	if bcrypt.CompareHashAndPassword([]byte(row.Password), []byte(password)) != nil {
		return domain.User{}, errors.New("неверное имя пользователя или пароль")
	}
	return domain.User{ID: row.ID, Username: row.Username, Role: row.Role}, nil
}

func CreateUser(db *gorm.DB, username, password string) error {
	if strings.TrimSpace(username) == "" || password == "" {
		return errors.New("логин и пароль обязательны")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	return db.Create(&userRow{Username: strings.TrimSpace(username), Password: string(hash), Role: "user"}).Error
}

func ListRecords(db *gorm.DB, userID uint) ([]domain.SZIRecord, error) {
	var rows []recordRow
	if err := db.Where("owner_id = ? OR (owner_id = 0 AND user_id = ?)", userID, userID).Order("name").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]domain.SZIRecord, len(rows))
	for i := range rows {
		result[i] = toDomain(rows[i])
	}
	return result, nil
}

func CreateRecord(db *gorm.DB, record domain.SZIRecord) error {
	return db.Create(fromDomain(record)).Error
}
func UpdateRecord(db *gorm.DB, record domain.SZIRecord) error {
	return db.Where("id = ? AND (owner_id = ? OR (owner_id = 0 AND user_id = ?))", record.ID, record.OwnerID, record.OwnerID).Updates(fromDomain(record)).Error
}
func DeleteRecord(db *gorm.DB, userID, recordID uint) error {
	return db.Where("id = ? AND (owner_id = ? OR (owner_id = 0 AND user_id = ?))", recordID, userID, userID).Delete(&recordRow{}).Error
}

func Statistics(db *gorm.DB, userID uint) ([]domain.SZIRecord, error) { return ListRecords(db, userID) }

func toDomain(row recordRow) domain.SZIRecord {
	return domain.SZIRecord{ID: row.ID, Name: row.Name, SZIType: row.SZIType, ProtectionClass: row.ProtectionClass, TrustLevel: row.TrustLevel, ACClass: row.ACClass, Vendor: row.Vendor, Version: row.Version, CertificateNumber: row.CertificateNumber, CertificationScheme: row.CertificationScheme, IssueDate: row.IssueDate, ExpiryDate: row.ExpiryDate, InstallLocation: row.InstallLocation, ResponsiblePerson: row.ResponsiblePerson, Notes: row.Notes, Status: row.Status, OwnerID: row.OwnerID, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt}
}
func fromDomain(record domain.SZIRecord) *recordRow {
	return &recordRow{ID: record.ID, Name: record.Name, SZIType: record.SZIType, ProtectionClass: record.ProtectionClass, TrustLevel: record.TrustLevel, ACClass: record.ACClass, Vendor: record.Vendor, Version: record.Version, CertificateNumber: record.CertificateNumber, CertificationScheme: record.CertificationScheme, IssueDate: record.IssueDate, ExpiryDate: record.ExpiryDate, InstallLocation: record.InstallLocation, ResponsiblePerson: record.ResponsiblePerson, Notes: record.Notes, Status: record.Status, OwnerID: record.OwnerID, UserID: record.OwnerID}
}
func CalculateStatus(expiry, now time.Time) string {
	if expiry.Before(now) {
		return domain.StatusExpired
	}
	if expiry.Before(now.AddDate(0, 0, 60)) {
		return domain.StatusAttention
	}
	return domain.StatusCurrent
}
func ValidateRecord(record domain.SZIRecord) error {
	if strings.TrimSpace(record.Name) == "" {
		return errors.New("наименование СЗИ обязательно")
	}
	if strings.TrimSpace(record.CertificateNumber) == "" {
		return errors.New("номер сертификата обязателен")
	}
	if record.IssueDate.IsZero() || record.ExpiryDate.IsZero() {
		return errors.New("даты сертификата обязательны")
	}
	if record.IssueDate.After(time.Now()) {
		return errors.New("дата выдачи не может быть в будущем")
	}
	if record.ExpiryDate.Before(record.IssueDate) {
		return errors.New("срок действия не может быть раньше даты выдачи")
	}
	return nil
}
