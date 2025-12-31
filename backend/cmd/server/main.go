package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/tau-tau-run/backend/config"
	"github.com/tau-tau-run/backend/internal/database"
	"github.com/tau-tau-run/backend/internal/middleware"
	"github.com/tau-tau-run/backend/internal/services"
	"github.com/tau-tau-run/backend/internal/utils"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	utils.ServerLogger.Info("Starting Tau-Tau Run API Server")
	utils.ServerLogger.Info("Environment: %s", cfg.Server.Env)

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize services
	authService := services.NewAuthService(cfg)

	// Setup Gin
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	
	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.CORS(cfg))
	router.Use(middleware.ErrorHandler())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		// Check database connection
		if err := database.HealthCheck(); err != nil {
			middleware.RespondWithError(c, 500, "UNHEALTHY", "Database connection failed", nil)
			return
		}

		middleware.RespondWithSuccess(c, 200, "", gin.H{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		public := v1.Group("/public")
		{
			public.GET("/health", func(c *gin.Context) {
				middleware.RespondWithSuccess(c, 200, "", gin.H{
					"status":  "healthy",
					"version": "1.0.0",
				})
			})
			// POST /register will be added in Phase 3
		}

		// Admin routes
		admin := v1.Group("/admin")
		{
			// POST /login (no auth required)
			// Will be added in Phase 4
			
			// Protected admin routes
			protected := admin.Group("")
			protected.Use(middleware.AuthMiddleware(authService))
			{
				// GET /participants will be added in Phase 5
				// PATCH /participants/:id/payment will be added in Phase 5
			}
		}
	}

	// Start server
	port := cfg.Server.Port
	utils.ServerLogger.Info("Server listening on port %s", port)
	utils.ServerLogger.Info("Public API: http://localhost:%s/api/v1/public", port)
	utils.ServerLogger.Info("Admin API: http://localhost:%s/api/v1/admin", port)

	// Graceful shutdown
	go func() {
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.ServerLogger.Info("Shutting down server...")
	fmt.Println("\n✅ Server shutdown complete")
}
