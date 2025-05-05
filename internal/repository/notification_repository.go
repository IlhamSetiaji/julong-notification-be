package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-notification-be/database"
	"github.com/IlhamSetiaji/julong-notification-be/internal/entity"
	"github.com/IlhamSetiaji/julong-notification-be/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type INotificationRepository interface {
	CreateNotification(ent *entity.Notification) (*entity.Notification, error)
	GetNotificationsByKeys(keys map[string]interface{}) ([]entity.Notification, error)
	GetNotificationsByKeysPagination(keys map[string]interface{}, page, pageSize int, search string, sort map[string]interface{}) ([]entity.Notification, int64, error)
	GetAllNotifications() ([]entity.Notification, error)
	FindByKeys(keys map[string]interface{}) (*entity.Notification, error)
	UpdateNotification(ent *entity.Notification) (*entity.Notification, error)
	DeleteNotification(id uuid.UUID) error
	GetUnreadNotificationCount(userID uuid.UUID, application string) (int64, error)
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
	if keys["read_at"] != nil {
		readAt := keys["read_at"]
		delete(keys, "read_at")
		if readAt == "YES" {
			query := r.db.GetDb().Where(keys).Where("read_at IS NOT NULL")
			r.log.GetLogger().Info("Generated SQL query ", "query", query.Debug().Statement.SQL.String())
			err := query.Find(&ent).Error
			if err != nil {
				r.log.GetLogger().Error("Failed to get notifications by keys: ", "error", err)
				return nil, err
			}
		} else if readAt == "NO" {
			err := r.db.GetDb().Where(keys).Where("read_at IS NULL").Find(&ent).Error
			if err != nil {
				r.log.GetLogger().Error("Failed to get notifications by keys: ", "error", err)
				return nil, err
			}
		} else {
			r.log.GetLogger().Error("Invalid value for read_at key: ", "value", readAt)
			return nil, errors.New("invalid value for read_at key")
		}
	} else {
		err := r.db.GetDb().Where(keys).Find(&ent).Error
		if err != nil {
			r.log.GetLogger().Error("Failed to get notifications by keys: ", "error", err)
			return nil, err
		}
	}
	return ent, nil
}

func (r *NotificationRepository) GetNotificationsByKeysPagination(keys map[string]interface{}, page, pageSize int, search string, sort map[string]interface{}) ([]entity.Notification, int64, error) {
	ent := []entity.Notification{}
	count := int64(0)
	query := r.db.GetDb().Model(&entity.Notification{}).Where(keys)

	if search != "" {
		query = query.Where("name ILIKE ? OR message ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if len(sort) > 0 {
		for column, order := range sort {
			query = query.Order(column + " " + order.(string))
		}
	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Count(&count).Error; err != nil {
		r.log.GetLogger().Error("Failed to count notifications: ", "error", err)
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&ent).Error; err != nil {
		r.log.GetLogger().Error("Failed to get notifications by keys with pagination: ", "error", err)
		return nil, 0, err
	}

	return ent, count, nil
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error
		}
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

func (r *NotificationRepository) GetUnreadNotificationCount(userID uuid.UUID, application string) (int64, error) {
	ent := int64(0)
	err := r.db.GetDb().Model(&entity.Notification{}).
		Where("user_id = ? AND application = ? AND read_at IS NULL", userID, application).
		Count(&ent).Error
	if err != nil {
		r.log.GetLogger().Error("Failed to get unread notification count: ", "error", err)
		return 0, err
	}
	return ent, nil
}
