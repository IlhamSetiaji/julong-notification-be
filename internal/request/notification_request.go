package request

type CreateNotificationRequest struct {
	Application string   `json:"application" validate:"required,application"`
	Name        string   `json:"name" validate:"required,name"`
	URL         string   `json:"url" validate:"required,url"`
	Message     string   `json:"message" validate:"required,message"`
	UserIDs     []string `json:"user_ids" validate:"required,dive"`
	CreatedBy   string   `json:"created_by" validate:"required,uuid"`
}

type UpdateNotificationRequest struct {
	ID          string `json:"id" validate:"required,uuid"`
	Application string `json:"application" validate:"omitempty,application"`
	Name        string `json:"name" validate:"omitempty,name"`
	URL         string `json:"url" validate:"omitempty,url"`
	Message     string `json:"message" validate:"omitempty,message"`
	CreatedBy   string `json:"created_by" validate:"omitempty,uuid"`
}
