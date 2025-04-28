package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-notification-be/internal/request"
	"github.com/IlhamSetiaji/julong-notification-be/internal/usecase"
	"github.com/IlhamSetiaji/julong-notification-be/logger"
	"github.com/IlhamSetiaji/julong-notification-be/utils"
	"github.com/IlhamSetiaji/julong-notification-be/validator"
	"github.com/gin-gonic/gin"
)

type INotificationHandler interface {
	CreateNotification(ctx *gin.Context)
	GetNotificationsByKeys(ctx *gin.Context)
	GetAllNotifications(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	GetByUserID(ctx *gin.Context)
	UpdateNotification(ctx *gin.Context)
	DeleteNotification(ctx *gin.Context)
}

type NotificationHandler struct {
	logger              logger.Logger
	validator           validator.Validator
	notificationUseCase usecase.INotificationUseCase
}

func NewNotificationHandler(
	logger logger.Logger,
	validator validator.Validator,
	notificationUseCase usecase.INotificationUseCase) INotificationHandler {
	return &NotificationHandler{
		logger:              logger,
		validator:           validator,
		notificationUseCase: notificationUseCase,
	}
}

func (h *NotificationHandler) CreateNotification(ctx *gin.Context) {
	var req request.CreateNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.GetLogger().Error("Failed to bind JSON: ", "error", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to bind JSON", err.Error())
		return
	}

	if err := h.validator.GetValidator().Struct(req); err != nil {
		h.logger.GetLogger().Error("Validation error: ", "error", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	err := h.notificationUseCase.CreateNotification(&req)
	if err != nil {
		h.logger.GetLogger().Error("Failed to create notification: ", "error", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create notification", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Notification created successfully", nil)
}

func (h *NotificationHandler) GetNotificationsByKeys(ctx *gin.Context) {
	application := ctx.Query("application")
	userID := ctx.Query("user_id")
	keys := make(map[string]interface{})
	if application != "" {
		keys["application"] = application
	}

	if userID != "" {
		keys["user_id"] = userID
	}

	notifications, err := h.notificationUseCase.GetNotificationsByKeys(keys)
	if err != nil {
		h.logger.GetLogger().Error("Failed to get notifications by keys: ", "error", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get notifications", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Notifications retrieved successfully", notifications)
}

func (h *NotificationHandler) GetAllNotifications(ctx *gin.Context) {
	notifications, err := h.notificationUseCase.GetAllNotifications()
	if err != nil {
		h.logger.GetLogger().Error("Failed to get all notifications: ", "error", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get notifications", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Notifications retrieved successfully", notifications)
}

func (h *NotificationHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	notification, err := h.notificationUseCase.FindByID(id)
	if err != nil {
		h.logger.GetLogger().Error("Failed to find notification by ID: ", "error", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find notification", err.Error())
		return
	}

	if notification == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Notification not found", "Notification not found")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Notification retrieved successfully", notification)
}

func (h *NotificationHandler) GetByUserID(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	notifications, err := h.notificationUseCase.GetByUserID(userID)
	if err != nil {
		h.logger.GetLogger().Error("Failed to get notifications by user ID: ", "error", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get notifications", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Notifications retrieved successfully", notifications)
}

func (h *NotificationHandler) UpdateNotification(ctx *gin.Context) {
	var req request.UpdateNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.GetLogger().Error("Failed to bind JSON: ", "error", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to bind JSON", err.Error())
		return
	}

	if err := h.validator.GetValidator().Struct(req); err != nil {
		h.logger.GetLogger().Error("Validation error: ", "error", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	res, err := h.notificationUseCase.UpdateNotification(&req)
	if err != nil {
		h.logger.GetLogger().Error("Failed to update notification: ", "error", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update notification", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Notification updated successfully", res)
}

func (h *NotificationHandler) DeleteNotification(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.notificationUseCase.DeleteNotification(id)
	if err != nil {
		h.logger.GetLogger().Error("Failed to delete notification: ", "error", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete notification", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Notification deleted successfully", nil)
}
