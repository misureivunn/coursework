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

	// Связь с пользователем
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
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

// AuditLog — история изменений для аудита
type AuditLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	User        *User     `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`
	RecordID    uint      `json:"record_id" gorm:"index"`                         // ID записи СЗИ (может быть NULL если запись удалена)
	Action      string    `json:"action" gorm:"size:50;not null"`                 // CREATE, UPDATE, DELETE
	EntityType  string    `json:"entity_type" gorm:"size:50;not null"`            // SZIRecord, User
	OldValue    string    `json:"old_value,omitempty" gorm:"type:text"`           // JSON со старыми значениями
	NewValue    string    `json:"new_value,omitempty" gorm:"type:text"`           // JSON с новыми значениями
	Description string    `json:"description,omitempty" gorm:"type:text"`
	IPAddress   string    `json:"ip_address,omitempty" gorm:"size:45"`
	CreatedAt   time.Time `json:"created_at" gorm:"index"`
}

// Notification — уведомления о событиях
type Notification struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"not null;index"`
	User      *User      `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	RecordID  uint       `json:"record_id" gorm:"index"` // ID записи СЗИ
	Record    *SZIRecord `json:"record,omitempty" gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE"`
	Type      string     `json:"type" gorm:"size:50;not null"`      // EXPIRING_SOON, EXPIRED, CRITICAL
	Title     string     `json:"title" gorm:"size:255;not null"`
	Message   string     `json:"message" gorm:"type:text;not null"`
	IsRead    bool       `json:"is_read" gorm:"default:false;index"`
	Priority  int        `json:"priority" gorm:"default:1"` // 1=низкий, 2=средний, 3=высокий, 4=критический
	CreatedAt time.Time  `json:"created_at" gorm:"index"`
}
