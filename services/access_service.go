package services

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"szi-registry/models"
	"szi-registry/repositories"
	"time"
)

// предоставляет доступ к записи другому пользователю
func GrantAccess(db *gorm.DB, ownerUserID, guestUserID, recordID uint, permission string, expiresAt *time.Time) error {
	permissionModel := &models.AccessPermission{
		OwnerUserID: ownerUserID,
		GuestUserID: guestUserID,
		RecordID:    recordID,
		Permission:  permission,
		ExpiresAt:   expiresAt,
	}

	repo := repositories.NewAccessPermissionRepository(db)
	err := repo.Create(permissionModel)
	if err != nil {
		return err
	}

	// Отправляем уведомление пользователю, которому предоставлен доступ
	notification := &models.Notification{
		SenderID:    ownerUserID,
		RecipientID: guestUserID,
		Title:       "Предоставлен доступ к записи СЗИ",
		Message:     fmt.Sprintf("Вам предоставлен доступ (%s) к записи СЗИ с ID %d", permission, recordID),
		Type:        "access_granted",
	}

	notificationRepo := repositories.NewNotificationRepository(db)
	err = notificationRepo.Create(notification)
	if err != nil {
		// Не прерываем операцию, если не удалось отправить уведомление
		log.Printf("Предупреждение: не удалось отправить уведомление: %v", err)
	}

	return nil
}

// отменяет доступ к записи другому пользователю
func RevokeAccess(db *gorm.DB, ownerUserID, guestUserID, recordID uint) error {
	repo := repositories.NewAccessPermissionRepository(db)
	err := repo.Delete(ownerUserID, guestUserID, recordID)
	if err != nil {
		return err
	}

	// Отправляем уведомление пользователю, у которого отозван доступ
	notification := &models.Notification{
		SenderID:    ownerUserID,
		RecipientID: guestUserID,
		Title:       "Доступ к записи СЗИ отозван",
		Message:     fmt.Sprintf("У вас отозван доступ к записи СЗИ с ID %d", recordID),
		Type:        "access_revoked",
	}

	notificationRepo := repositories.NewNotificationRepository(db)
	err = notificationRepo.Create(notification)
	if err != nil {
		// Не прерываем операцию, если не удалось отправить уведомление
		log.Printf("Предупреждение: не удалось отправить уведомление об отзыве доступа: %v", err)
	}

	return nil
}

// возвращает пользователя по имени
func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	repo := repositories.NewUserRepository(db)
	return repo.FindByUsername(username)
}

// возвращает записи, доступные пользователю
func GetUserRecords(db *gorm.DB, userID uint) ([]models.SZIRecord, error) {
	service := &SziService{DB: db}
	return service.GetSziRecordsAccessibleByUserID(userID)
}

// удаляет запись СЗИ
func DeleteSziRecord(db *gorm.DB, userID, recordID uint) error {
	repo := repositories.NewSZIRecordRepository(db)
	return repo.DeleteForUser(userID, recordID)
}

// возвращает всех пользователей системы
func GetAllUsers(db *gorm.DB) ([]models.User, error) {
	repo := repositories.NewUserRepository(db)
	return repo.GetAll()
}

// возвращает разрешения, предоставленные пользователю
func GetUserPermissions(db *gorm.DB, userID uint) ([]models.AccessPermission, error) {
	repo := repositories.NewAccessPermissionRepository(db)
	return repo.FindByGuestUserID(userID)
}

// возвращает разрешения, предоставленные пользователем
func GetUserGrantedPermissions(db *gorm.DB, userID uint) ([]models.AccessPermission, error) {
	repo := repositories.NewAccessPermissionRepository(db)
	return repo.FindByOwnerUserID(userID)
}

// назначает роль пользователю
func AssignRole(db *gorm.DB, userID uint, roleName, description string) error {
	role := &models.UserRole{
		UserID:      userID,
		RoleName:    roleName,
		Description: description,
	}

	repo := repositories.NewUserRoleRepository(db)
	return repo.Create(role)
}

// удаляет роль у пользователя
func RemoveRole(db *gorm.DB, userID uint, roleName string) error {
	repo := repositories.NewUserRoleRepository(db)
	return repo.Delete(userID, roleName)
}

// возвращает роли пользователя
func GetUserRoles(db *gorm.DB, userID uint) ([]models.UserRole, error) {
	repo := repositories.NewUserRoleRepository(db)
	return repo.FindByUserID(userID)
}

// возвращает всех пользователей с определенной ролью
func GetAllUsersWithRole(db *gorm.DB, roleName string) ([]models.User, error) {
	repo := repositories.NewUserRoleRepository(db)
	return repo.FindUsersWithRole(roleName)
}

// отправляет уведомление пользователю
func SendNotification(db *gorm.DB, senderID, recipientID uint, title, message, notificationType string) error {
	notification := &models.Notification{
		SenderID:    senderID,
		RecipientID: recipientID,
		Title:       title,
		Message:     message,
		Type:        notificationType,
	}

	repo := repositories.NewNotificationRepository(db)
	return repo.Create(notification)
}

// возвращает уведомления пользователя
func GetNotifications(db *gorm.DB, userID uint) ([]models.Notification, error) {
	repo := repositories.NewNotificationRepository(db)
	return repo.FindByRecipientID(userID)
}

// помечает уведомление как прочитанное
func MarkNotificationAsRead(db *gorm.DB, notificationID uint) error {
	repo := repositories.NewNotificationRepository(db)
	return repo.MarkAsRead(notificationID)
}

// отправляет запрос на доступ к записи
func RequestAccess(db *gorm.DB, requesterID, ownerID, recordID uint, permission, message string) error {
	request := &models.AccessRequest{
		RequesterID: requesterID,
		OwnerID:     ownerID,
		RecordID:    recordID,
		Permission:  permission,
		Message:     message,
		Status:      "pending",
	}

	repo := repositories.NewAccessRequestRepository(db)
	return repo.Create(request)
}

// одобряет запрос на доступ
func ApproveAccessRequest(db *gorm.DB, requestID, processorID uint) error {
	repo := repositories.NewAccessRequestRepository(db)
	return repo.Approve(requestID, processorID)
}

// отклоняет запрос на доступ
func RejectAccessRequest(db *gorm.DB, requestID, processorID uint, reason string) error {
	repo := repositories.NewAccessRequestRepository(db)
	return repo.Reject(requestID, processorID, reason)
}

// возвращает ожидающие запросы на доступ
func GetPendingAccessRequests(db *gorm.DB, ownerID uint) ([]models.AccessRequest, error) {
	repo := repositories.NewAccessRequestRepository(db)
	return repo.FindPendingByOwnerID(ownerID)
}

// создает новую запись СЗИ
func CreateSZIRecord(db *gorm.DB, userID uint, record *models.SZIRecord) (uint, error) {
	record.UserID = userID

	repo := repositories.NewSZIRecordRepository(db)
	err := repo.Create(record)
	if err != nil {
		return 0, err
	}

	return record.ID, nil
}

// возвращает конкретную запись СЗИ
func GetSZIRecord(db *gorm.DB, userID, recordID uint) (*models.SZIRecord, error) {
	repo := repositories.NewSZIRecordRepository(db)
	record, err := repo.FindByID(recordID)
	if err != nil {
		return nil, err
	}

	// Проверяем, принадлежит ли запись пользователю или ему предоставлен доступ
	if record.UserID != userID {
		// Проверяем, есть ли у пользователя доступ к этой записи
		accessRepo := repositories.NewAccessPermissionRepository(db)
		hasAccess, err := accessRepo.HasAccess(userID, recordID)
		if err != nil || !hasAccess {
			return nil, gorm.ErrRecordNotFound
		}
	}

	return record, nil
}

// обновляет запись СЗИ
func UpdateSZIRecord(db *gorm.DB, userID, recordID uint, record *models.SZIRecord) error {
	// Проверяем, принадлежит ли запись пользователю
	repo := repositories.NewSZIRecordRepository(db)
	existingRecord, err := repo.FindByID(recordID)
	if err != nil {
		return err
	}

	if existingRecord.UserID != userID {
		return gorm.ErrRecordNotFound
	}

	record.ID = recordID
	return repo.Update(record)
}

// удаляет запись СЗИ
func DeleteSZIRecord(db *gorm.DB, userID, recordID uint) error {
	// Проверяем, принадлежит ли запись пользователю
	repo := repositories.NewSZIRecordRepository(db)
	existingRecord, err := repo.FindByID(recordID)
	if err != nil {
		return err
	}

	if existingRecord.UserID != userID {
		return gorm.ErrRecordNotFound
	}

	return repo.Delete(recordID)
}

// ищет записи СЗИ по критериям
func SearchSZIRecords(db *gorm.DB, userID uint, criteria map[string]interface{}) ([]models.SZIRecord, error) {
	service := &SziService{DB: db}
	return service.SearchSziRecords(userID, criteria)
}

// сервис для работы с записями СЗИ
type SziService struct {
	DB *gorm.DB
}

// возвращает все записи СЗИ, доступные пользователю
func (s *SziService) GetSziRecordsAccessibleByUserID(userID uint) ([]models.SZIRecord, error) {
	repo := repositories.NewSZIRecordRepository(s.DB)
	return repo.FindAccessibleByUserID(userID)
}

// выполняет поиск записей СЗИ по различным критериям
func (s *SziService) SearchSziRecords(userID uint, filters map[string]interface{}) ([]models.SZIRecord, error) {
	repo := repositories.NewSZIRecordRepository(s.DB)
	return repo.Search(userID, filters)
}
