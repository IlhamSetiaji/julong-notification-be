// websocket_handler.go
package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-notification-be/internal/websocket"
	"github.com/IlhamSetiaji/julong-notification-be/logger"
	"github.com/IlhamSetiaji/julong-notification-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IWebSocketHandler interface {
	HandleWebSocket(c *gin.Context)
}

type WebSocketHandler struct {
	log logger.Logger
	hub *websocket.Hub
}

func NewWebSocketHandler(log logger.Logger, hub *websocket.Hub) IWebSocketHandler {
	return &WebSocketHandler{
		log: log,
		hub: hub,
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	userIDStr := c.Query("user_id")
	appType := c.Query("app_type")

	if userIDStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "User ID is required", "User ID is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID format", "Invalid user ID format")
		return
	}

	websocket.ServeWS(h.hub, c.Writer, c.Request, userID, appType)
}
