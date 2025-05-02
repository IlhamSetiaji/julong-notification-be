package usecase

import (
	"errors"
	"time"

	"github.com/IlhamSetiaji/julong-notification-be/internal/dto"
	"github.com/IlhamSetiaji/julong-notification-be/internal/entity"
	"github.com/IlhamSetiaji/julong-notification-be/internal/repository"
	"github.com/IlhamSetiaji/julong-notification-be/internal/request"
	"github.com/IlhamSetiaji/julong-notification-be/internal/response"
	"github.com/IlhamSetiaji/julong-notification-be/internal/websocket"
	"github.com/IlhamSetiaji/julong-notification-be/logger"
	"github.com/google/uuid"
)

type INotificationUseCase interface {
	CreateNotification(req *request.CreateNotificationRequest) error
	GetNotificationsByKeys(keys map[string]interface{}) ([]response.NotificationResponse, error)
	GetAllNotifications() ([]response.NotificationResponse, error)
	FindByID(id string) (*response.NotificationResponse, error)
	GetByUserID(userID string) ([]response.NotificationResponse, error)
	UpdateNotification(req *request.UpdateNotificationRequest) (*response.NotificationResponse, error)
	DeleteNotification(id string) error
	GetUnreadNotificationCount(userID string, application string) (int64, error)
}

type NotificationUseCase struct {
	log                    logger.Logger
	notificationDTO        dto.INotificationDTO
	notificationRepository repository.INotificationRepository
	hub                    *websocket.Hub
}

func NewNotificationUseCase(
	log logger.Logger,
	notificationDTO dto.INotificationDTO,
	notificationRepository repository.INotificationRepository,
	hub *websocket.Hub) INotificationUseCase {
	return &NotificationUseCase{
		log:                    log,
		notificationDTO:        notificationDTO,
		notificationRepository: notificationRepository,
		hub:                    hub,
	}
}

func (uc *NotificationUseCase) CreateNotification(req *request.CreateNotificationRequest) error {
	if len(req.UserIDs) == 0 {
		return errors.New("user_ids cannot be empty")
	}
	createdByUUID, err := uuid.Parse(req.CreatedBy)
	if err != nil {
		return errors.New("invalid created_by format")
	}

	for _, userID := range req.UserIDs {
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return errors.New("invalid user_id format")
		}

		notification := &entity.Notification{
			Application: req.Application,
			Name:        req.Name,
			URL:         req.URL,
			Message:     req.Message,
			UserID:      userUUID,
			CreatedBy:   createdByUUID,
		}

		createdNotification, err := uc.notificationRepository.CreateNotification(notification)
		if err != nil {
			uc.log.GetLogger().Error("Failed to create notification: ", err)
			return err
		}

		unreadCount, err := uc.notificationRepository.GetUnreadNotificationCount(userUUID, req.Application)
		if err != nil {
			uc.log.GetLogger().Error("Failed to get unread notification count: ", err)
			return err
		}

		createdNotification.UnreadCount = unreadCount

		wsNotification := uc.notificationDTO.ConvertEntityToWebsocketResponse(createdNotification)
		uc.hub.BroadcastNotification(*wsNotification)
	}

	return nil
}

func (uc *NotificationUseCase) GetNotificationsByKeys(keys map[string]interface{}) ([]response.NotificationResponse, error) {
	notifications, err := uc.notificationRepository.GetNotificationsByKeys(keys)
	if err != nil {
		uc.log.GetLogger().Error("Failed to get notifications by keys: ", err)
		return nil, err
	}

	var responses []response.NotificationResponse
	for _, notification := range notifications {
		response := uc.notificationDTO.ConvertEntityToResponse(&notification)
		responses = append(responses, *response)
	}

	return responses, nil
}

func (uc *NotificationUseCase) GetAllNotifications() ([]response.NotificationResponse, error) {
	notifications, err := uc.notificationRepository.GetAllNotifications()
	if err != nil {
		uc.log.GetLogger().Error("Failed to get all notifications: ", err)
		return nil, err
	}

	var responses []response.NotificationResponse
	for _, notification := range notifications {
		response := uc.notificationDTO.ConvertEntityToResponse(&notification)
		responses = append(responses, *response)
	}

	return responses, nil
}

func (uc *NotificationUseCase) FindByID(id string) (*response.NotificationResponse, error) {
	notification, err := uc.notificationRepository.FindByKeys(map[string]interface{}{"id": id})
	if err != nil {
		uc.log.GetLogger().Error("Failed to find notification by ID: ", err)
		return nil, err
	}

	if notification == nil {
		return nil, nil
	}

	response := uc.notificationDTO.ConvertEntityToResponse(notification)
	return response, nil
}

func (uc *NotificationUseCase) GetByUserID(userID string) ([]response.NotificationResponse, error) {
	notifications, err := uc.notificationRepository.GetNotificationsByKeys(map[string]interface{}{"user_id": userID})
	if err != nil {
		uc.log.GetLogger().Error("Failed to get notifications by user ID: ", err)
		return nil, err
	}

	var responses []response.NotificationResponse
	for _, notification := range notifications {
		response := uc.notificationDTO.ConvertEntityToResponse(&notification)
		responses = append(responses, *response)
	}

	return responses, nil
}

func (uc *NotificationUseCase) UpdateNotification(req *request.UpdateNotificationRequest) (*response.NotificationResponse, error) {
	notification, err := uc.notificationRepository.FindByKeys(map[string]interface{}{"id": req.ID})
	if err != nil {
		uc.log.GetLogger().Error("Failed to find notification by ID: ", err)
		return nil, err
	}

	if notification == nil {
		return nil, errors.New("notification not found")
	}

	if req.Application != "" {
		notification.Application = req.Application
	}
	if req.Name != "" {
		notification.Name = req.Name
	}
	if req.URL != "" {
		notification.URL = req.URL
	}
	if req.Message != "" {
		notification.Message = req.Message
	}
	if req.ReadAt != nil {
		// parsedTime, err := time.Parse("2006-01-02 15:04:05", *req.ReadAt)
		parsedTime, err := time.Parse(time.RFC3339, *req.ReadAt)
		if err != nil {
			return nil, errors.New("invalid ReadAt format, must be RFC3339")
		}
		notification.ReadAt = &parsedTime
	}

	_, err = uc.notificationRepository.UpdateNotification(notification)
	if err != nil {
		return nil, err
	}

	response := uc.notificationDTO.ConvertEntityToResponse(notification)
	return response, nil
}

func (uc *NotificationUseCase) DeleteNotification(id string) error {
	notification, err := uc.notificationRepository.FindByKeys(map[string]interface{}{"id": id})
	if err != nil {
		uc.log.GetLogger().Error("Failed to find notification by ID: ", err)
		return err
	}

	if notification == nil {
		return errors.New("notification not found")
	}

	err = uc.notificationRepository.DeleteNotification(notification.ID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *NotificationUseCase) GetUnreadNotificationCount(userID string, application string) (int64, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		uc.log.GetLogger().Error("Invalid user ID format: ", err)
		return 0, errors.New("invalid user ID format")
	}
	notification, err := uc.notificationRepository.GetUnreadNotificationCount(parsedUserID, application)
	if err != nil {
		uc.log.GetLogger().Error("Failed to get unread notification count: ", err)
		return 0, err
	}
	return notification, nil
}
