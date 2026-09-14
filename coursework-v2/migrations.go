package main

import (
	"szi-registry/models"

	"gorm.io/gorm"
)

func normalizeLegacyReferenceFields(db *gorm.DB) error {
	var records []models.SZIRecord
	if err := db.Where("szi_type_reference = '' OR szi_type_reference IS NULL").Find(&records).Error; err != nil {
		return err
	}
	for index := range records {
		record := &records[index]
		switch record.Type {
		case "СКЗИ":
			record.SZIType = "Средство криптографической защиты информации"
		case "ОС":
			record.SZIType = "Операционная система"
		case "СЗИ КС", "СЗИ СКЗИ":
			record.SZIType = "Средство защиты от НСД"
		default:
			record.SZIType = "Другое"
		}
		record.ProtectionClass = record.ClassProtection
		record.Vendor = record.Manufacturer
		record.Version = record.SoftwareVersion
		record.CertificateNumber = record.CertNumber
		record.IssueDate = record.CertIssueDate
		record.ExpiryDate = record.CertExpiryDate
		record.InstallLocation = record.Location
		if err := db.Save(record).Error; err != nil {
			return err
		}
	}
	return nil
}
