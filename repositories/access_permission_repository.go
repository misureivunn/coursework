package repositories

import (
	"errors"
	"gorm.io/gorm"
	"szi-registry/models"
)

type AccessPermissionRepository struct {
	DB *gorm.DB
}

func NewAccessPermissionRepository(db *gorm.DB) *AccessPermissionRepository {
	return &AccessPermissionRepository{DB: db}
}

func (r *AccessPermissionRepository) Create(permission *models.AccessPermission) error {
	// Проверяем, что владелец записи совпадает с владельцем, дающим доступ
	var record models.SZIRecord
	err := r.DB.Where("id = ?", permission.RecordID).First(&record).Error
	if err != nil {
		return err
	}

	if record.UserID != permission.OwnerUserID {
		return errors.New("только владелец записи может предоставить доступ к ней")
	}

	// Проверяем, существует ли уже такое разрешение
	var existingPermission models.AccessPermission
	err = r.DB.Where("owner_user_id = ? AND guest_user_id = ? AND record_id = ?",
		permission.OwnerUserID, permission.GuestUserID, permission.RecordID).First(&existingPermission).Error

	if err == nil {
		// Обновляем существующее разрешение
		existingPermission.Permission = permission.Permission
		existingPermission.ExpiresAt = permission.ExpiresAt
		return r.DB.Save(&existingPermission).Error
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Создаем новое разрешение
	return r.DB.Create(permission).Error
}

func (r *AccessPermissionRepository) Delete(ownerUserID, guestUserID, recordID uint) error {
	result := r.DB.Where("owner_user_id = ? AND guest_user_id = ? AND record_id = ?",
		ownerUserID, guestUserID, recordID).Delete(&models.AccessPermission{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("разрешение не найдено или уже удалено")
	}

	return nil
}

func (r *AccessPermissionRepository) FindByGuestUserID(userID uint) ([]models.AccessPermission, error) {
	var permissions []models.AccessPermission
	err := r.DB.Where("guest_user_id = ?", userID).Order("created_at DESC").Find(&permissions).Error
	return permissions, err
}

func (r *AccessPermissionRepository) FindByOwnerUserID(userID uint) ([]models.AccessPermission, error) {
	var permissions []models.AccessPermission
	err := r.DB.Where("owner_user_id = ?", userID).Order("created_at DESC").Find(&permissions).Error
	return permissions, err
}

func (r *AccessPermissionRepository) HasAccess(userID, recordID uint) (bool, error) {
	var permission models.AccessPermission
	err := r.DB.Where("guest_user_id = ? AND record_id = ? AND (expires_at IS NULL OR expires_at > NOW())",
		userID, recordID).First(&permission).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
