package response

import (
	"time"

	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID            uuid.UUID  `json:"id"`
	Application   string     `json:"application"`
	Name          string     `json:"name"`
	URL           string     `json:"url"`
	ReadAt        *time.Time `json:"read_at"`
	Message       string     `json:"message"`
	UserID        uuid.UUID  `json:"user_id"`
	CreatedBy     uuid.UUID  `json:"created_by"`
	UserName      string     `json:"user_name"`
	CreatedByName string     `json:"created_by_name"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
