package request

type NotificationRequest struct {
	Application string   `json:"application" validate:"required,application"`
	Name        string   `json:"name" validate:"required,name"`
	URL         string   `json:"url" validate:"required,url"`
	Message     string   `json:"message" validate:"required,message"`
	UserIDs     []string `json:"user_ids" validate:"required,dive"`
	CreatedBy   string   `json:"created_by" validate:"required,uuid"`
}
