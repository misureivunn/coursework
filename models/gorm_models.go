package models

import (
	"time"

	"gorm.io/gorm"
)

// основная модель
type SZIRecord struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"size:255;not null;index:idx_szi_user_name,priority:2"`
	Type            string    `json:"szi_type" gorm:"column:szi_type;size:50;not null;index:idx_szi_type"`
	CertNumber      string    `json:"cert_number" gorm:"column:cert_number;size:64;not null;index:idx_szi_user_cert,priority:2"`
	CertIssueDate   time.Time `json:"cert_issue_date" gorm:"column:cert_issue_date;not null;index:idx_szi_issue_date"`
	CertExpiryDate  time.Time `json:"cert_expiry_date" gorm:"column:cert_expiry_date;not null;index:idx_szi_expiry_date"`
	Location        string    `json:"location" gorm:"size:255;not null"`
	Status          string    `json:"status" gorm:"size:20;not null;default:Актуально;index:idx_szi_status"`
	UserID          uint      `json:"user_id" gorm:"column:user_id;not null;index:idx_szi_user,priority:1"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Manufacturer    string    `json:"manufacturer,omitempty" gorm:"size:255"`
	SoftwareVersion string    `json:"software_version,omitempty" gorm:"column:software_version;size:64"`

	// Поля для классификации СЗИ от НСД
	Purpose         string `json:"purpose,omitempty" gorm:"column:purpose;size:50;index:idx_szi_purpose"`
	DeploymentType  string `json:"deployment_type,omitempty" gorm:"column:deployment_type;size:50;index:idx_szi_deploy"`
	ClassProtection string `json:"class_protection,omitempty" gorm:"column:class_protection;size:10;index:idx_szi_class"`
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
	Username  string          `json:"username" gorm:"size:64;unique;not null;index:idx_user_username"`
	Password  string          `json:"password" gorm:"size:255;not null"` // хешированный пароль
	Role      string          `json:"role" gorm:"size:32;default:user;index:idx_user_role"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
