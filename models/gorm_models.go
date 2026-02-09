package models

import (
	"gorm.io/gorm"
	"time"
)

// основная мрдель
type SZIRecord struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name" gorm:"not null"`
	Type              string    `json:"szi_type" gorm:"column:szi_type;not null"`
	CertNumber        string    `json:"cert_number" gorm:"column:cert_number;not null"`
	CertIssueDate     time.Time `json:"cert_issue_date" gorm:"column:cert_issue_date;not null"`
	CertExpiryDate    time.Time `json:"cert_expiry_date" gorm:"column:cert_expiry_date;not null"`
	Location          string    `json:"location" gorm:"not null"`
	Status            string    `json:"status" gorm:"default:Актуально"`
	UserID            uint      `json:"user_id" gorm:"column:user_id;not null"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Manufacturer      string    `json:"manufacturer,omitempty"`
	SoftwareVersion   string    `json:"software_version,omitempty" gorm:"column:software_version"`

	// Поля для классификации СЗИ от НСД
	Purpose         string `json:"purpose,omitempty" gorm:"column:purpose"`                         // Назначение (АС, ИВК, Универсальное)
	DeploymentType  string `json:"deployment_type,omitempty" gorm:"column:deployment_type"`         // Тип развертывания (Клиент-сервер, Автономное, АПК, Виртуальное)
	ClassProtection string `json:"class_protection,omitempty" gorm:"column:class_protection"`      // Класс защищенности (1, 2, 3А, 4, 5)
}


type UserRole struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"column:user_id;not null"`
	RoleName    string    `json:"role_name" gorm:"column:role_name;not null"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}


type User struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	Username  string          `json:"username" gorm:"unique;not null"`
	Password  string          `json:"password" gorm:"not null"` // хешированный пароль
	Role      string          `json:"role" gorm:"default:user"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

