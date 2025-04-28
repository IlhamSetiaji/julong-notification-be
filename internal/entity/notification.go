package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model  `json:"-"`
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	Application string     `json:"application" gorm:"type:varchar(255);not null"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	URL         string     `json:"url" gorm:"type:text;not null"`
	ReadAt      *time.Time `json:"read_at" gorm:"type:timestamp"`
	Message     string     `json:"message" gorm:"type:text;not null"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	CreatedBy   uuid.UUID  `json:"created_by" gorm:"type:uuid;not null"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	n.ID = uuid.New()
	loc := time.FixedZone("Asia/Jakarta", 7*60*60)
	n.CreatedAt = time.Now().In(loc)
	n.UpdatedAt = time.Now().In(loc)
	return
}

func (n *Notification) BeforeUpdate(tx *gorm.DB) (err error) {
	loc := time.FixedZone("Asia/Jakarta", 7*60*60)
	n.UpdatedAt = time.Now().In(loc)
	return
}

func (Notification) TableName() string {
	return "notifications"
}
