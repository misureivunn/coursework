package domain

import "time"

type User struct {
	ID       uint
	Username string
	Password string
	Role     string
}

type SZIRecord struct {
	ID                  uint
	Name                string
	SZIType             string
	ProtectionClass     string
	TrustLevel          string
	ACClass             string
	Vendor              string
	Version             string
	CertificateNumber   string
	CertificationScheme string
	IssueDate           time.Time
	ExpiryDate          time.Time
	InstallLocation     string
	ResponsiblePerson   string
	Notes               string
	Status              string
	OwnerID             uint
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

const (
	StatusCurrent   = "Актуально"
	StatusAttention = "Требует внимания"
	StatusExpired   = "Просрочено"
)

var SZITypeValues = []string{
	"Средство защиты от НСД",
	"Средство криптографической защиты информации",
	"Операционная система",
	"Межсетевой экран",
	"Система обнаружения вторжений",
	"Средство антивирусной защиты",
	"Средство доверенной загрузки",
	"Средство контроля съёмных носителей",
	"Другое",
}

var ProtectionClassValues = []string{"1", "2", "3", "4", "5", "6"}
var TrustLevelValues = []string{"1", "2", "3", "4", "5", "6"}
var ACClassValues = []string{"1Д", "1Г", "1В", "1Б", "1А", "2Б", "2А", "3Б", "3А", "не применимо"}
var CertificationSchemeValues = []string{"Серия", "Единичный образец"}

func StatusFor(expiry, now time.Time) string {
	if expiry.Before(now) {
		return StatusExpired
	}
	if expiry.Before(now.AddDate(0, 0, 60)) {
		return StatusAttention
	}
	return StatusCurrent
}
