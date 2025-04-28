package main

import (
	"github.com/IlhamSetiaji/julong-notification-be/config"
	"github.com/IlhamSetiaji/julong-notification-be/database"
	"github.com/IlhamSetiaji/julong-notification-be/logger"
	"github.com/IlhamSetiaji/julong-notification-be/server"
	"github.com/IlhamSetiaji/julong-notification-be/validator"
)

func main() {
	config := config.GetConfig()
	logger := logger.NewLogger()
	db := database.NewPostgresDatabase(config)
	validator := validator.NewValidatorV10(config)
	server := server.NewGinServer(db, *config, logger, validator)

	// Start the server
	server.Start()
}
