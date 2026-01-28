package repositories

import (
	"errors"
	"gorm.io/gorm"
	"szi-registry/models"
)

type UserRoleRepository struct {
	DB *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{DB: db}
}

func (r *UserRoleRepository) Create(role *models.UserRole) error {
	// Проверяем, существует ли уже такая роль у пользователя
	var existingRole models.UserRole
	err := r.DB.Where("user_id = ? AND role_name = ?", role.UserID, role.RoleName).First(&existingRole).Error

	if err == nil {
		// Обновляем существующую роль
		existingRole.Description = role.Description
		return r.DB.Save(&existingRole).Error
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Создаем новую роль
	return r.DB.Create(role).Error
}

func (r *UserRoleRepository) Delete(userID uint, roleName string) error {
	result := r.DB.Where("user_id = ? AND role_name = ?", userID, roleName).Delete(&models.UserRole{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("роль не найдена или уже удалена")
	}

	return nil
}

func (r *UserRoleRepository) FindByUserID(userID uint) ([]models.UserRole, error) {
	var roles []models.UserRole
	err := r.DB.Where("user_id = ?", userID).Find(&roles).Error
	return roles, err
}

func (r *UserRoleRepository) FindUsersWithRole(roleName string) ([]models.User, error) {
	var users []models.User
	err := r.DB.Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_name = ?", roleName).
		Find(&users).Error
	return users, err
}
