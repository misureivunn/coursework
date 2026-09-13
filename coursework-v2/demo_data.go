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
			Name:            "Dallas Lock 8.0-К",
			Type:            "СЗИ КС",
			CertNumber:      "Сведения уточняются по реестру ФСТЭК",
			CertIssueDate:   now.AddDate(-2, 0, 0),
			CertExpiryDate:  now.AddDate(1, 0, 0),
			Location:        "Серверная",
			Manufacturer:    "Конфидент",
			SoftwareVersion: "8.0-К",
			Purpose:         "АС",
			DeploymentType:  "Клиент-сервер",
			ClassProtection: "1",
			Notes:           "Демонстрационная запись по существующему продукту. Реквизиты сертификата необходимо сверить в открытом реестре ФСТЭК.",
			UserID:          userID,
		},
		{
			Name:            "Secret Net Studio 8",
			Type:            "СЗИ КС",
			CertNumber:      "Сведения уточняются по реестру ФСТЭК",
			CertIssueDate:   now.AddDate(-1, 6, 0),
			CertExpiryDate:  now.AddDate(0, 8, 0),
			Location:        "Рабочие станции",
			Manufacturer:    "ИнфоТеКС",
			SoftwareVersion: "8",
			Purpose:         "АС",
			DeploymentType:  "Клиент-сервер",
			ClassProtection: "2",
			Notes:           "Демонстрационная запись по существующему продукту. Реквизиты сертификата необходимо сверить в открытом реестре ФСТЭК.",
			UserID:          userID,
		},
		{
			Name:            "Аккорд-Win32",
			Type:            "СКЗИ",
			CertNumber:      "Сведения уточняются по реестру ФСТЭК",
			CertIssueDate:   now.AddDate(-3, 0, 0),
			CertExpiryDate:  now.AddDate(0, 1, 15),
			Location:        "АРМ администратора",
			Manufacturer:    "Аккорд",
			SoftwareVersion: "Win32",
			Purpose:         "ИВК",
			DeploymentType:  "Автономное",
			ClassProtection: "3А",
			Notes:           "Демонстрационная запись по существующему продукту. Реквизиты сертификата необходимо сверить в открытом реестре ФСТЭК.",
			UserID:          userID,
		},
		{
			Name:            "Astra Linux Special Edition",
			Type:            "ОС",
			CertNumber:      "Сведения уточняются по реестру ФСТЭК",
			CertIssueDate:   now.AddDate(-1, 0, 0),
			CertExpiryDate:  now.AddDate(-0, -1, 0),
			Location:        "Виртуальная инфраструктура",
			Manufacturer:    "РусБИТех",
			SoftwareVersion: "1.7",
			Purpose:         "Универсальное",
			DeploymentType:  "Виртуальное",
			ClassProtection: "4",
			Notes:           "Демонстрационная запись по существующему продукту. Срок выбран для показа статуса «Просрочено». Реквизиты сертификата необходимо сверить в открытом реестре ФСТЭК.",
			UserID:          userID,
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
