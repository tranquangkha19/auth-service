package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tranquangkha19/auth-service/internal/auth"
	"github.com/tranquangkha19/auth-service/internal/config"
	"github.com/tranquangkha19/auth-service/internal/database"
)

type Server struct {
	router *gin.Engine
	cfg    config.Config
	db     *database.Database
}

func NewServer(cfg config.Config) *Server {
	router := gin.Default()

	// Initialize database
	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	authHandler := auth.NewHandler(auth.NewService(db))

	router.POST("/login", authHandler.Login)
	router.POST("/register", authHandler.Register)
	router.POST("/validate-token", authHandler.ValidateToken)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": cfg.AppName,
		})
	})

	return &Server{
		router: router,
		cfg:    cfg,
		db:     db,
	}
}

func (s *Server) Run() {
	s.router.Run(fmt.Sprintf(":%s", s.cfg.Port))
}
