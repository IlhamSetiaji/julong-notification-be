package dto

import (
	"github.com/IlhamSetiaji/julong-notification-be/internal/entity"
	"github.com/IlhamSetiaji/julong-notification-be/internal/messaging"
	"github.com/IlhamSetiaji/julong-notification-be/internal/request"
	"github.com/IlhamSetiaji/julong-notification-be/internal/response"
	"github.com/IlhamSetiaji/julong-notification-be/internal/websocket"
	"github.com/IlhamSetiaji/julong-notification-be/logger"
)

type INotificationDTO interface {
	ConvertEntityToResponse(ent *entity.Notification) *response.NotificationResponse
	ConvertEntityToWebsocketResponse(ent *entity.Notification) *websocket.WsNotification
}

type NotificationDTO struct {
	log         logger.Logger
	userMessage messaging.IUserMessage
}

func NewNotificationDTO(log logger.Logger, userMessage messaging.IUserMessage) INotificationDTO {
	return &NotificationDTO{
		log:         log,
		userMessage: userMessage,
	}
}

func (n *NotificationDTO) ConvertEntityToResponse(ent *entity.Notification) *response.NotificationResponse {
	var userName string
	var createdByName string
	user, err := n.userMessage.SendFindUserByIDMessage(request.SendFindUserByIDMessageRequest{
		ID: ent.UserID.String(),
	})
	if err != nil {
		n.log.GetLogger().Error("Failed to find user by ID: ", "error", err)
		userName = "Unknown"
	} else {
		userName = user.Name
	}

	createdBy, err := n.userMessage.SendFindUserByIDMessage(request.SendFindUserByIDMessageRequest{
		ID: ent.CreatedBy.String(),
	})
	if err != nil {
		n.log.GetLogger().Error("Failed to find user by ID: ", "error", err)
		createdByName = "Unknown"
	} else {
		createdByName = createdBy.Name
	}

	return &response.NotificationResponse{
		ID:            ent.ID,
		Application:   ent.Application,
		Name:          ent.Name,
		URL:           ent.URL,
		ReadAt:        ent.ReadAt,
		Message:       ent.Message,
		UserID:        ent.UserID,
		CreatedBy:     ent.CreatedBy,
		UserName:      userName,
		CreatedByName: createdByName,
		CreatedAt:     ent.CreatedAt,
		UpdatedAt:     ent.UpdatedAt,
	}
}

func (n *NotificationDTO) ConvertEntityToWebsocketResponse(ent *entity.Notification) *websocket.WsNotification {
	var userName string
	var createdByName string
	user, err := n.userMessage.SendFindUserByIDMessage(request.SendFindUserByIDMessageRequest{
		ID: ent.UserID.String(),
	})
	if err != nil {
		n.log.GetLogger().Error("Failed to find user by ID: ", "error", err)
		userName = "Unknown"
	} else {
		userName = user.Name
	}

	createdBy, err := n.userMessage.SendFindUserByIDMessage(request.SendFindUserByIDMessageRequest{
		ID: ent.CreatedBy.String(),
	})
	if err != nil {
		n.log.GetLogger().Error("Failed to find user by ID: ", "error", err)
		createdByName = "Unknown"
	} else {
		createdByName = createdBy.Name
	}

	return &websocket.WsNotification{
		ID:            ent.ID,
		Application:   ent.Application,
		Name:          ent.Name,
		URL:           ent.URL,
		ReadAt:        ent.ReadAt,
		Message:       ent.Message,
		UserID:        ent.UserID,
		CreatedBy:     ent.CreatedBy,
		UserName:      userName,
		CreatedByName: createdByName,
		CreatedAt:     ent.CreatedAt,
		UpdatedAt:     ent.UpdatedAt,
	}
}
