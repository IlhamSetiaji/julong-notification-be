package server

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/IlhamSetiaji/julong-notification-be/config"
	"github.com/IlhamSetiaji/julong-notification-be/database"
	"github.com/IlhamSetiaji/julong-notification-be/internal/rabbitmq"
	"github.com/IlhamSetiaji/julong-notification-be/logger"
	"github.com/IlhamSetiaji/julong-notification-be/validator"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type ginServer struct {
	app       *gin.Engine
	db        database.Database
	conf      config.Config
	log       logger.Logger
	validator validator.Validator
}

func NewGinServer(db database.Database, conf config.Config, log logger.Logger, validator validator.Validator) Server {
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(gin.Logger())

	store := cookie.NewStore([]byte(conf.Session.Secret))
	app.Use(sessions.Sessions(conf.Session.Name, store))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	app.Use(func(c *gin.Context) {
		if !shouldExcludeFromCSRF(c.Request.URL.Path) {
			csrf.Middleware(csrf.Options{
				Secret: conf.Csrf.Secret,
				ErrorFunc: func(c *gin.Context) {
					c.String(http.StatusForbidden, "CSRF token mismatch")
					c.Abort()
				},
			})(c)
		}
		c.Next()
	})

	// app.RedirectTrailingSlash = false

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		rabbitmq.InitConsumer(conf, log)
	}()

	go func() {
		defer wg.Done()
		rabbitmq.InitProducer(conf, log)
	}()

	return &ginServer{
		app:       app,
		db:        db,
		conf:      conf,
		log:       log,
		validator: validator,
	}
}

func (g *ginServer) Start() {
	g.app.Static("/storage", "./storage")
	g.app.Static("/assets", "./public")

	g.app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("App-Name", g.conf.Server.Name)
	})
	g.app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Service is running",
			"status":  "OK",
		})
	})

	g.log.GetLogger().Info("Server started on port " + strconv.Itoa(g.conf.Server.Port))
	g.app.Run(":" + strconv.Itoa(g.conf.Server.Port))
}

func (g *ginServer) GetApp() *gin.Engine {
	return g.app
}

func shouldExcludeFromCSRF(path string) bool {
	return len(path) >= 4 && path[:4] == "/api"
}
