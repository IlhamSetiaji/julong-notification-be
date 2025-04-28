package request

type CreateNotificationRequest struct {
	Application string   `json:"application" validate:"required,application"`
	Name        string   `json:"name" validate:"required"`
	URL         string   `json:"url" validate:"required"`
	Message     string   `json:"message" validate:"required"`
	UserIDs     []string `json:"user_ids" validate:"required,dive"`
	CreatedBy   string   `json:"created_by" validate:"required,uuid"`
}

type UpdateNotificationRequest struct {
	ID          string  `json:"id" validate:"required,uuid"`
	Application string  `json:"application" validate:"omitempty,application"`
	Name        string  `json:"name" validate:"omitempty"`
	URL         string  `json:"url" validate:"omitempty"`
	Message     string  `json:"message" validate:"omitempty"`
	ReadAt      *string `json:"read_at" validate:"omitempty"`
	CreatedBy   string  `json:"created_by" validate:"omitempty,uuid"`
}
