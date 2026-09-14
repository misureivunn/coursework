package main

import (
	"time"

	"szi-registry/models"
	"szi-registry/services"

	"gorm.io/gorm"
)

// seedDemoRecords добавляет демонстрационные записи только в пустой реестр пользователя.
func seedDemoRecords(db *gorm.DB, userID uint) error {
	var count int64
	if err := db.Model(&models.SZIRecord{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	now := time.Now()
	records := []models.SZIRecord{
		{
			Name:                "Dallas Lock 8.0-К",
			Type:                "СЗИ КС",
			SZIType:             "Средство защиты от НСД",
			ProtectionClass:     "3",
			TrustLevel:          "6",
			ACClass:             "1Б",
			CertNumber:          "Сведения уточняются по реестру ФСТЭК",
			CertIssueDate:       now.AddDate(-2, 0, 0),
			CertExpiryDate:      now.AddDate(1, 0, 0),
			Location:            "Серверная",
			Manufacturer:        "Конфидент",
			Vendor:              "Конфидент",
			SoftwareVersion:     "8.0-К",
			Version:             "8.0-К",
			CertificateNumber:   "3745",
			CertificationScheme: "Серия",
			IssueDate:           now.AddDate(-2, 0, 0),
			ExpiryDate:          now.AddDate(1, 0, 0),
			InstallLocation:     "Серверная",
			ResponsiblePerson:   "Администратор ИБ",
			Purpose:             "АС",
			DeploymentType:      "Клиент-сервер",
			ClassProtection:     "1",
			Notes:               "Демонстрационная запись по существующему продукту. Реквизиты сертификата необходимо сверить в открытом реестре ФСТЭК.",
			UserID:              userID,
		},
		{
			Name:                "Secret Net Studio 8",
			Type:                "СЗИ КС",
			SZIType:             "Средство защиты от НСД",
			ProtectionClass:     "5",
			TrustLevel:          "4",
			ACClass:             "1В",
			CertNumber:          "Сведения уточняются по реестру ФСТЭК",
			CertIssueDate:       now.AddDate(-1, 6, 0),
			CertExpiryDate:      now.AddDate(0, 8, 0),
			Location:            "Рабочие станции",
			Manufacturer:        "ИнфоТеКС",
			Vendor:              "ИнфоТеКС",
			SoftwareVersion:     "8",
			Version:             "8.4.2863.0",
			CertificateNumber:   "246/7",
			CertificationScheme: "Единичный образец",
			IssueDate:           now.AddDate(-1, 6, 0),
			ExpiryDate:          now.AddDate(0, 8, 0),
			InstallLocation:     "Рабочие станции",
			ResponsiblePerson:   "Специалист ИБ",
			Purpose:             "АС",
			DeploymentType:      "Клиент-сервер",
			ClassProtection:     "2",
			Notes:               "Демонстрационная запись по существующему продукту. Реквизиты сертификата необходимо сверить в открытом реестре ФСТЭК.",
			UserID:              userID,
		},
		{
			Name:                "Аккорд-Win32",
			Type:                "СКЗИ",
			SZIType:             "Средство криптографической защиты информации",
			ProtectionClass:     "",
			TrustLevel:          "5",
			ACClass:             "3Б",
			CertNumber:          "Сведения уточняются по реестру ФСТЭК",
			CertIssueDate:       now.AddDate(-3, 0, 0),
			CertExpiryDate:      now.AddDate(0, 1, 15),
			Location:            "АРМ администратора",
			Manufacturer:        "Аккорд",
			Vendor:              "Аккорд",
			SoftwareVersion:     "Win32",
			Version:             "Win32",
			CertificateNumber:   "Сведения уточняются по реестру ФСТЭК",
			CertificationScheme: "Серия",
			IssueDate:           now.AddDate(-3, 0, 0),
			ExpiryDate:          now.AddDate(0, 1, 15),
			InstallLocation:     "АРМ администратора",
			ResponsiblePerson:   "Администратор ИБ",
			Purpose:             "ИВК",
			DeploymentType:      "Автономное",
			ClassProtection:     "3А",
			Notes:               "Демонстрационная запись по существующему продукту. Реквизиты сертификата необходимо сверить в открытом реестре ФСТЭК.",
			UserID:              userID,
		},
		{
			Name:                "Astra Linux Special Edition",
			Type:                "ОС",
			SZIType:             "Операционная система",
			ProtectionClass:     "4",
			TrustLevel:          "4",
			ACClass:             "2А",
			CertNumber:          "Сведения уточняются по реестру ФСТЭК",
			CertIssueDate:       now.AddDate(-1, 0, 0),
			CertExpiryDate:      now.AddDate(-0, -1, 0),
			Location:            "Виртуальная инфраструктура",
			Manufacturer:        "РусБИТех",
			Vendor:              "РусБИТех",
			SoftwareVersion:     "1.7",
			Version:             "1.7",
			CertificateNumber:   "Сведения уточняются по реестру ФСТЭК",
			CertificationScheme: "Серия",
			IssueDate:           now.AddDate(-1, 0, 0),
			ExpiryDate:          now.AddDate(0, -1, 0),
			InstallLocation:     "Виртуальная инфраструктура",
			ResponsiblePerson:   "Архитектор ИБ",
			Purpose:             "Универсальное",
			DeploymentType:      "Виртуальное",
			ClassProtection:     "4",
			Notes:               "Демонстрационная запись по существующему продукту. Срок выбран для показа статуса «Просрочено». Реквизиты сертификата необходимо сверить в открытом реестре ФСТЭК.",
			UserID:              userID,
		},
	}

	for index := range records {
		if err := services.CreateSziRecord(db, &records[index]); err != nil {
			return err
		}
	}
	return nil
}

func seedDemoRecordsForExistingUsers(db *gorm.DB) error {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}
	for _, user := range users {
		if err := seedDemoRecords(db, user.ID); err != nil {
			return err
		}
	}
	return nil
}

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
