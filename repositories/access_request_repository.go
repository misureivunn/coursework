package repositories

import (
	"errors"
	"gorm.io/gorm"
	"szi-registry/models"
)

type AccessRequestRepository struct {
	DB *gorm.DB
}

func NewAccessRequestRepository(db *gorm.DB) *AccessRequestRepository {
	return &AccessRequestRepository{DB: db}
}

func (r *AccessRequestRepository) Create(request *models.AccessRequest) error {
	// Проверяем, что владелец записи совпадает с владельцем, к которому запрашивается доступ
	var record models.SZIRecord
	err := r.DB.Where("id = ?", request.RecordID).First(&record).Error
	if err != nil {
		return err
	}

	if record.UserID != request.OwnerID {
		return errors.New("указанный владелец не является владельцем записи")
	}

	return r.DB.Create(request).Error
}

func (r *AccessRequestRepository) Approve(requestID, processorID uint) error {
	// Получаем информацию о запросе
	request := models.AccessRequest{}
	err := r.DB.First(&request, requestID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("запрос доступа не найден или уже обработан")
		}
		return err
	}

	// Проверяем, что обработчик является владельцем записи
	if request.OwnerID != processorID {
		return errors.New("только владелец записи может одобрить запрос доступа")
	}

	// Обновляем статус запроса
	request.Status = "approved"
	request.ProcessedAt = &request.UpdatedAt
	request.ProcessedBy = &processorID

	err = r.DB.Save(&request).Error
	if err != nil {
		return err
	}

	// Предоставляем доступ
	permission := &models.AccessPermission{
		OwnerUserID: request.OwnerID,
		GuestUserID: request.RequesterID,
		RecordID:    request.RecordID,
		Permission:  request.Permission,
	}

	permissionRepo := NewAccessPermissionRepository(r.DB)
	err = permissionRepo.Create(permission)
	if err != nil {
		return err
	}

	return nil
}

func (r *AccessRequestRepository) Reject(requestID, processorID uint, reason string) error {
	// Получаем информацию о запросе
	request := models.AccessRequest{}
	err := r.DB.First(&request, requestID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("запрос доступа не найден или уже обработан")
		}
		return err
	}

	// Проверяем, что обработчик является владельцем записи
	if request.OwnerID != processorID {
		return errors.New("только владелец записи может отклонить запрос доступа")
	}

	// Обновляем статус запроса
	request.Status = "rejected"
	request.ProcessedAt = &request.UpdatedAt
	request.ProcessedBy = &processorID

	return r.DB.Save(&request).Error
}

func (r *AccessRequestRepository) FindPendingByOwnerID(ownerID uint) ([]models.AccessRequest, error) {
	var requests []models.AccessRequest
	err := r.DB.Where("owner_id = ? AND status = 'pending'", ownerID).Order("created_at ASC").Find(&requests).Error
	return requests, err
}
