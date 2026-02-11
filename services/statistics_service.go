package services

import (
	"szi-registry/models"
	"time"
	"gorm.io/gorm"
)

// Структура для хранения статистики по СЗИ
type SZIStatistics struct {
	TotalCount          int                    `json:"total_count"`
	CertificateStats    CertificateStatistics  `json:"certificate_stats"`
	ProtectionClassStats map[string]int        `json:"protection_class_stats"`
	TypeStats           map[string]int         `json:"type_stats"`
	PurposeStats        map[string]int         `json:"purpose_stats"`
	DeploymentTypeStats map[string]int         `json:"deployment_type_stats"`
}

// Структура для статистики по сертификатам
type CertificateStatistics struct {
	Active     int `json:"active"`
	Expired    int `json:"expired"`
	NoCert     int `json:"no_cert"`
}

// Функция для получения статистики по СЗИ
func GetSZIStatistics(db *gorm.DB, userID uint) (*SZIStatistics, error) {
	var records []models.SZIRecord
	err := db.Where("user_id = ?", userID).Find(&records).Error
	if err != nil {
		return nil, err
	}

	stats := &SZIStatistics{
		CertificateStats: CertificateStatistics{},
		ProtectionClassStats: make(map[string]int),
		TypeStats: make(map[string]int),
		PurposeStats: make(map[string]int),
		DeploymentTypeStats: make(map[string]int),
	}

	now := time.Now()
	
	for _, record := range records {
		stats.TotalCount++
		
		// Статистика по сертификатам
		if record.CertExpiryDate.IsZero() {
			stats.CertificateStats.NoCert++
		} else if record.CertExpiryDate.Before(now) {
			stats.CertificateStats.Expired++
		} else {
			stats.CertificateStats.Active++
		}
		
		// Статистика по классам защиты
		if record.ClassProtection != "" {
			stats.ProtectionClassStats[record.ClassProtection]++
		}
		
		// Статистика по типам СЗИ
		if record.Type != "" {
			stats.TypeStats[record.Type]++
		}
		
		// Статистика по назначению
		if record.Purpose != "" {
			stats.PurposeStats[record.Purpose]++
		}
		
		// Статистика по типу развертывания
		if record.DeploymentType != "" {
			stats.DeploymentTypeStats[record.DeploymentType]++
		}
	}

	return stats, nil
}