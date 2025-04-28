package repository

import (
	"github.com/IlhamSetiaji/julong-notification-be/database"
	"github.com/IlhamSetiaji/julong-notification-be/internal/entity"
	"github.com/IlhamSetiaji/julong-notification-be/logger"
	"github.com/google/uuid"
)

type INotificationRepository interface {
	CreateNotification(ent *entity.Notification) (*entity.Notification, error)
	GetNotificationsByKeys(keys map[string]interface{}) ([]entity.Notification, error)
	GetAllNotifications() ([]entity.Notification, error)
	FindByKeys(keys map[string]interface{}) (*entity.Notification, error)
	UpdateNotification(ent *entity.Notification) (*entity.Notification, error)
	DeleteNotification(id uuid.UUID) error
}

type NotificationRepository struct {
	db  database.Database
	log logger.Logger
}

func NewNotificationRepository(db database.Database, log logger.Logger) INotificationRepository {
	return &NotificationRepository{
		db:  db,
		log: log,
	}
}

func (r *NotificationRepository) CreateNotification(ent *entity.Notification) (*entity.Notification, error) {
	err := r.db.GetDb().Create(ent).Error
	if err != nil {
		r.log.GetLogger().Error("Failed to create notification: ", "error", err)
		return nil, err
	}
	return ent, nil
}

func (r *NotificationRepository) GetNotificationsByKeys(keys map[string]interface{}) ([]entity.Notification, error) {
	ent := []entity.Notification{}
	err := r.db.GetDb().Where(keys).Find(&ent).Error
	if err != nil {
		r.log.GetLogger().Error("Failed to get notifications by keys: ", "error", err)
		return nil, err
	}
	return ent, nil
}

func (r *NotificationRepository) GetAllNotifications() ([]entity.Notification, error) {
	ent := []entity.Notification{}
	err := r.db.GetDb().Find(&ent).Error
	if err != nil {
		r.log.GetLogger().Error("Failed to get all notifications: ", "error", err)
		return nil, err
	}
	return ent, nil
}

func (r *NotificationRepository) FindByKeys(keys map[string]interface{}) (*entity.Notification, error) {
	ent := &entity.Notification{}
	err := r.db.GetDb().Where(keys).First(ent).Error
	if err != nil {
		r.log.GetLogger().Error("Failed to find notification by keys: ", "error", err)
		return nil, err
	}

	return ent, nil
}

func (r *NotificationRepository) UpdateNotification(ent *entity.Notification) (*entity.Notification, error) {
	err := r.db.GetDb().Where("id = ?", ent.ID).Updates(ent).Error
	if err != nil {
		r.log.GetLogger().Error("Failed to update notification: ", "error", err)
		return nil, err
	}
	return ent, nil
}

func (r *NotificationRepository) DeleteNotification(id uuid.UUID) error {
	err := r.db.GetDb().Where("id = ?", id).Delete(&entity.Notification{}).Error
	if err != nil {
		r.log.GetLogger().Error("Failed to delete notification: ", "error", err)
		return err
	}
	return nil
}
