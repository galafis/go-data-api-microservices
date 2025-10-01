package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/galafis/go-data-api-microservices/internal/config"
	"github.com/galafis/go-data-api-microservices/internal/handlers"
	"github.com/galafis/go-data-api-microservices/internal/middleware"
	"github.com/galafis/go-data-api-microservices/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Data API Microservices
// @version         1.0
// @description     A high-performance data API microservices system
// @termsOfService  http://swagger.io/terms/

// @contact.name   Gabriel Demetrios Lafis
// @contact.url    https://github.com/galafis
// @contact.email  gabriel.lafis@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found, using environment variables")
	}

	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", err)
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(logger.GinLogger())
	router.Use(middleware.RequestID())
	router.Use(middleware.Metrics())

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.CORS.AllowOrigins
	corsConfig.AllowMethods = cfg.CORS.AllowMethods
	corsConfig.AllowHeaders = cfg.CORS.AllowHeaders
	corsConfig.ExposeHeaders = cfg.CORS.ExposeHeaders
	corsConfig.AllowCredentials = cfg.CORS.AllowCredentials
	corsConfig.MaxAge = cfg.CORS.MaxAge
	router.Use(cors.New(corsConfig))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.POST("/refresh", handlers.RefreshToken)
			auth.POST("/logout", middleware.AuthRequired(), handlers.Logout)
		}

		// Data routes
		data := v1.Group("/data")
		data.Use(middleware.AuthRequired())
		{
			data.GET("/datasets", handlers.ListDatasets)
			data.POST("/datasets", handlers.CreateDataset)
			data.GET("/datasets/:id", handlers.GetDataset)
			data.PUT("/datasets/:id", handlers.UpdateDataset)
			data.DELETE("/datasets/:id", handlers.DeleteDataset)
			
			data.POST("/query", handlers.QueryData)
			data.POST("/transform", handlers.TransformData)
			data.POST("/aggregate", handlers.AggregateData)
			data.POST("/join", handlers.JoinData)
		}

		// Analytics routes
		analytics := v1.Group("/analytics")
		analytics.Use(middleware.AuthRequired())
		{
			analytics.GET("/summary", handlers.GetDataSummary)
			analytics.POST("/statistics", handlers.ComputeStatistics)
			analytics.POST("/correlation", handlers.ComputeCorrelation)
			analytics.POST("/timeseries", handlers.AnalyzeTimeSeries)
			analytics.POST("/forecast", handlers.GenerateForecast)
		}

		// User routes
		users := v1.Group("/users")
		users.Use(middleware.AuthRequired())
		{
			users.GET("/me", handlers.GetCurrentUser)
			users.PUT("/me", handlers.UpdateCurrentUser)
			users.DELETE("/me", handlers.DeleteCurrentUser)
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		logger.Info(fmt.Sprintf("Server starting on port %d", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", err)
	}

	logger.Info("Server exited properly")
}

