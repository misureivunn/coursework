package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Password  string    `json:"password" gorm:"not null"`
	Role      string    `json:"role" gorm:"default:user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SZI struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name" gorm:"not null"`
	Type              string    `json:"type" gorm:"not null"`
	CertNumber        string    `json:"cert_number" gorm:"not null"`
	CertIssueDate     time.Time `json:"cert_issue_date"`
	CertExpiryDate    time.Time `json:"cert_expiry_date"`
	Location          string    `json:"location"`
	Status            string    `json:"status" gorm:"default:Актуально"`
	UserID            uint      `json:"user_id" gorm:"not null"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Manufacturer      string    `json:"manufacturer"`
	SoftwareVersion   string    `json:"software_version"`
	ContactPerson     string    `json:"contact_person"`
	DocumentationLink string    `json:"documentation_link"`
}