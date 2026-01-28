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
	ContactPerson     string    `json:"contact_person,omitempty" gorm:"column:contact_person"`
	DocumentationLink string    `json:"documentation_link,omitempty" gorm:"column:documentation_link"`
}

// права доступа
type AccessPermission struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	OwnerUserID uint       `json:"owner_user_id" gorm:"column:owner_user_id;not null"` // ID пользователя, которому принадлежит запись
	GuestUserID uint       `json:"guest_user_id" gorm:"column:guest_user_id;not null"` // ID пользователя, которому предоставлен доступ
	RecordID    uint       `json:"record_id" gorm:"column:record_id;not null"`         // ID записи СЗИ
	Permission  string     `json:"permission" gorm:"not null"`                         // r (read), rw (read-write)
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" gorm:"column:expires_at"` // время истечения доступа
}

type UserRole struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"column:user_id;not null"`
	RoleName    string    `json:"role_name" gorm:"column:role_name;not null"` // admin, viewer, editor и т.д.
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type Notification struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	SenderID    uint       `json:"sender_id" gorm:"column:sender_id;not null"`       // ID отправителя
	RecipientID uint       `json:"recipient_id" gorm:"column:recipient_id;not null"` // ID получателя
	Title       string     `json:"title" gorm:"not null"`                            // Заголовок уведомления
	Message     string     `json:"message" gorm:"not null"`                          // Текст уведомления
	Type        string     `json:"type" gorm:"not null"`                             // Тип уведомления: access_request, info, alert и т.д.
	ReadStatus  bool       `json:"read_status" gorm:"default:false"`                 // Статус прочтения
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ReadAt      *time.Time `json:"read_at,omitempty" gorm:"column:read_at"` // Время прочтения (опционально)
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

// запросы на доступ к записям
type AccessRequest struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	RequesterID uint       `json:"requester_id" gorm:"column:requester_id;not null"` // ID пользователя, запрашивающего доступ
	OwnerID     uint       `json:"owner_id" gorm:"column:owner_id;not null"`         // ID владельца записи
	RecordID    uint       `json:"record_id" gorm:"column:record_id;not null"`       // ID записи, к которой запрашивается доступ
	Permission  string     `json:"permission" gorm:"not null"`                       // Запрашиваемое разрешение: r, rw
	Status      string     `json:"status" gorm:"default:pending"`                    // pending, approved, rejected
	Message     string     `json:"message" gorm:"not null"`                          // Сообщение от запрашивающего
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty" gorm:"column:processed_at"` // Время обработки запроса
	ProcessedBy *uint      `json:"processed_by,omitempty" gorm:"column:processed_by"` // ID пользователя, обработавшего запрос
}
