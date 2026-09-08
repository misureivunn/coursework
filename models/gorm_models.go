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
	Notes           string    `json:"notes,omitempty" gorm:"type:text"`
	SoftwareVersion string    `json:"software_version,omitempty" gorm:"column:software_version;size:64"`
	Purpose         string    `json:"purpose,omitempty" gorm:"column:purpose;size:50;index:idx_szi_purpose"`
	DeploymentType  string    `json:"deployment_type,omitempty" gorm:"column:deployment_type;size:50;index:idx_szi_deploy"`
	ClassProtection string    `json:"class_protection,omitempty" gorm:"column:class_protection;size:10;index:idx_szi_class"`

	// Связь с пользователем
	User User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type User struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	Username  string          `json:"username" gorm:"size:64;unique;not null;index:idx_user_username"`
	Password  string          `json:"password" gorm:"size:255;not null"` // хешированный пароль
	Role      string          `json:"role" gorm:"size:32;default:user;index:idx_user_role"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Обратная связь с записями СЗИ
	SZIRecords []SZIRecord `json:"szi_records,omitempty" gorm:"foreignKey:UserID"`
}

// ReadNotification - модель для хранения прочитанных уведомлений
type ReadNotification struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"user_id" gorm:"index;not null"`
	NotificationKey string    `json:"notification_key" gorm:"not null"` // Уникальный ключ уведомления (например, "cert_expired_record_123")
	CreatedAt       time.Time `json:"created_at"`
}

// SZIAuditLog - модель для хранения истории изменений записей СЗИ
type SZIAuditLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	SZIID       uint      `json:"szi_id" gorm:"column:szi_id;not null;index:idx_audit_szi"`
	UserID      uint      `json:"user_id" gorm:"column:user_id;not null;index:idx_audit_user"`
	Action      string    `json:"action" gorm:"size:20;not null"` // create, update, delete
	FieldName   string    `json:"field_name,omitempty" gorm:"column:field_name;size:64"`
	OldValue    string    `json:"old_value,omitempty" gorm:"column:old_value;type:text"`
	NewValue    string    `json:"new_value,omitempty" gorm:"column:new_value;type:text"`
	PerformedAt time.Time `json:"performed_at" gorm:"column:performed_at;not null;index:idx_audit_date"`
	IPAddress   string    `json:"ip_address,omitempty" gorm:"column:ip_address;size:45"`

	// Связи
	SZI  SZIRecord `json:"szi,omitempty" gorm:"foreignKey:SZIID;constraint:OnDelete:CASCADE"`
	User User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
