package main

import (
	"github.com/IlhamSetiaji/julong-notification-be/config"
	"github.com/IlhamSetiaji/julong-notification-be/database"
	"github.com/IlhamSetiaji/julong-notification-be/internal/entity"
	"github.com/IlhamSetiaji/julong-notification-be/logger"
)

func main() {
	config := config.GetConfig()
	logger := logger.NewLogger()
	db := database.NewPostgresDatabase(config)

	// Initialize the database connection
	if err := db.GetDb().AutoMigrate(
		&entity.Notification{},
	); err != nil {
		logger.GetLogger().Fatal("Failed to migrate database", err)
	}

	logger.GetLogger().Info("Database migrated successfully")
}
