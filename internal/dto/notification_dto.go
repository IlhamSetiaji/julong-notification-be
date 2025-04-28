package dto

import (
	"github.com/IlhamSetiaji/julong-notification-be/internal/entity"
	"github.com/IlhamSetiaji/julong-notification-be/internal/response"
)

type INotificationDTO interface {
	ConvertEntityToResponse(ent *entity.Notification) *response.NotificationResponse
}

type NotificationDTO struct{}

func NewNotificationDTO() INotificationDTO {
	return &NotificationDTO{}
}

func (n *NotificationDTO) ConvertEntityToResponse(ent *entity.Notification) *response.NotificationResponse {
	return &response.NotificationResponse{
		ID:            ent.ID,
		Application:   ent.Application,
		Name:          ent.Name,
		URL:           ent.URL,
		ReadAt:        ent.ReadAt,
		Message:       ent.Message,
		UserID:        ent.UserID,
		CreatedBy:     ent.CreatedBy,
		UserName:      ent.UserName,
		CreatedByName: ent.CreatedByName,
		CreatedAt:     ent.CreatedAt,
		UpdatedAt:     ent.UpdatedAt,
	}
}
