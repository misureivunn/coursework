package main

import (
	"time"

	"szi-registry/models"
	"szi-registry/services"

	"gorm.io/gorm"
)

// seedDemoRecords добавляет демонстрационные записи только в пустой реестр пользователя.
func seedDemoRecords(db *gorm.DB, userID uint) error {
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
		{
			Name: "Аккорд-АМДЗ", SZIType: "Средство доверенной загрузки", ProtectionClass: "2", ACClass: "не применимо",
			Type: "Средство доверенной загрузки", Vendor: "АО «ОКБ САПР»", CertificateNumber: "246/7", CertificationScheme: "Серия",
			IssueDate: date(2007, 8, 13), ExpiryDate: date(2019, 8, 13), InstallLocation: "Архивный сервер (выведен из эксплуатации)", ResponsiblePerson: "Сидорова А.В.",
			CertNumber: "246/7", CertIssueDate: date(2007, 8, 13), CertExpiryDate: date(2019, 8, 13), Location: "Архивный сервер (выведен из эксплуатации)", Manufacturer: "АО «ОКБ САПР»",
			Notes: "Устаревший сертификат, приведён для демонстрации статуса «Просрочено».", UserID: userID,
		},
		{
			Name: "Astra Linux Special Edition РУСБ.10015-01 (обновление 1.8)", SZIType: "Операционная система", ProtectionClass: "1", TrustLevel: "1", ACClass: "1Б",
			Type: "Операционная система", Vendor: "ГК «Астра» (АО «НПО РусБИТех»)", CertificateNumber: "2557", CertificationScheme: "Серия", Version: "1.8",
			IssueDate: date(2012, 1, 27), ExpiryDate: date(2031, 1, 27), InstallLocation: "АРМ администратора ИБ", ResponsiblePerson: "Иванов И.И., администратор ИБ",
			CertNumber: "2557", CertIssueDate: date(2012, 1, 27), CertExpiryDate: date(2031, 1, 27), Location: "АРМ администратора ИБ", Manufacturer: "ГК «Астра» (АО «НПО РусБИТех»)", SoftwareVersion: "1.8",
			Notes: "Сертификат неоднократно переоформлялся при выходе новых обновлений (1.4-1.8).", UserID: userID,
		},
		{
			Name: "Kaspersky Endpoint Security для Windows (12.3.0.493)", SZIType: "Средство антивирусной защиты", ACClass: "1Г", Vendor: "АО «Лаборатория Касперского»", Version: "12.3.0.493",
			Type: "Средство антивирусной защиты", CertificateNumber: "4068", CertificationScheme: "Серия", IssueDate: date(2024, 3, 1), ExpiryDate: date(2029, 1, 22),
			InstallLocation: "АРМ отдела продаж", ResponsiblePerson: "Петров П.П.", CertNumber: "4068", CertIssueDate: date(2024, 3, 1), CertExpiryDate: date(2029, 1, 22), Location: "АРМ отдела продаж", Manufacturer: "АО «Лаборатория Касперского»", SoftwareVersion: "12.3.0.493",
			Notes: "Класс защиты САВЗ: Б2, В2, Г2 (приказ ФСТЭК №28) — не входит в шкалу СВТ 1-6.", UserID: userID,
		},
		{
			Name: "Kaspersky Endpoint Security для Linux (12.0.0.6672)", SZIType: "Средство антивирусной защиты", ACClass: "1Г", Vendor: "АО «Лаборатория Касперского»", Version: "12.0.0.6672",
			Type: "Средство антивирусной защиты", CertificateNumber: "2534", CertificationScheme: "Серия", IssueDate: date(2020, 12, 27), ExpiryDate: date(2030, 12, 27),
			InstallLocation: "Сервер приложений (Linux)", ResponsiblePerson: "Иванов И.И., администратор ИБ", CertNumber: "2534", CertIssueDate: date(2020, 12, 27), CertExpiryDate: date(2030, 12, 27), Location: "Сервер приложений (Linux)", Manufacturer: "АО «Лаборатория Касперского»", SoftwareVersion: "12.0.0.6672",
			Notes: "Дата выдачи оценочная — источник указывает только срок окончания действия.", UserID: userID,
		},
		{
			Name: "Kaspersky Endpoint Security для Linux ARM Edition (12.2.0.2412)", SZIType: "Средство антивирусной защиты", ACClass: "не применимо", Vendor: "АО «Лаборатория Касперского»", Version: "12.2.0.2412",
			Type: "Средство антивирусной защиты", CertificateNumber: "4672", CertificationScheme: "Серия", IssueDate: date(2024, 4, 1), ExpiryDate: date(2028, 5, 3),
			InstallLocation: "Сервер на архитектуре ARM (тестовый стенд)", ResponsiblePerson: "Сидорова А.В.", CertNumber: "4672", CertIssueDate: date(2024, 4, 1), CertExpiryDate: date(2028, 5, 3), Location: "Сервер на архитектуре ARM (тестовый стенд)", Manufacturer: "АО «Лаборатория Касперского»", SoftwareVersion: "12.2.0.2412",
			Notes: "Класс защиты САВЗ: Б2, В2, Г2, СКН П2 (приказ ФСТЭК №28) — не входит в шкалу СВТ 1-6.", UserID: userID,
		},
		{
			Name: "Kaspersky Endpoint Security для Mac (11.1.0.210a)", SZIType: "Средство антивирусной защиты", ACClass: "3Б", Vendor: "АО «Лаборатория Касперского»", Version: "11.1.0.210a",
			Type: "Средство антивирусной защиты", CertificateNumber: "3424", CertificationScheme: "Серия", IssueDate: date(2021, 7, 7), ExpiryDate: date(2026, 7, 7),
			InstallLocation: "АРМ дизайнера (macOS)", ResponsiblePerson: "Петров П.П.", CertNumber: "3424", CertIssueDate: date(2021, 7, 7), CertExpiryDate: date(2026, 7, 7), Location: "АРМ дизайнера (macOS)", Manufacturer: "АО «Лаборатория Касперского»", SoftwareVersion: "11.1.0.210a",
			Notes: "Срок действия истёк — демонстрация статуса «Просрочено» на актуальных данных.", UserID: userID,
		},
	}

	for index := range records {
		var existing models.SZIRecord
		if err := db.Where("user_id = ? AND name = ?", userID, records[index].Name).First(&existing).Error; err == nil {
			continue
		} else if err != gorm.ErrRecordNotFound {
			return err
		}
		if err := services.CreateSziRecord(db, &records[index]); err != nil {
			return err
		}
	}
	return nil
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
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
