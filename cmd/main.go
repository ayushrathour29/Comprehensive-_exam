package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"final_assessment/internal/db"
	"final_assessment/internal/handlers"
	"final_assessment/internal/services"
	"final_assessment/internal/utils"
)

func main() {
	// Initialize Logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Warn(utils.LogNoEnvFileFound)
	}
	fmt.Println("DATABASE_URL:", os.Getenv("DATABASE_URL"))

	// Connect to the database
	if err := db.Connect(); err != nil {
		logrus.WithError(err).Fatal(utils.LogDBConnectFail)
	}
	// Run migrations
	if err := db.RunMigrations(db.GetSQLDB()); err != nil {
		logrus.WithError(err).Fatal("Failed to run migrations")
	}
	defer db.Close()

	// Initialize Gin router
	r := gin.Default()

	// Create JobService and JobHandler
	jobService := services.NewJobService(db.DB)
	workerPool := services.NewWorkerPool(jobService, 5, 100)
	workerPool.Start()
	jobHandler := handlers.NewJobHandler(jobService, workerPool)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Job endpoints
	// POST /jobs -> Submit a new job
	// GET /jobs/:id -> Get job status and result
	// GET /jobs -> List jobs
	r.POST("/jobs", jobHandler.SubmitJob)
	r.GET("/jobs/:id", jobHandler.GetJob)
	r.GET("/jobs", jobHandler.ListJobs)

	// Start server
	port := os.Getenv(utils.EnvPort)
	if port == "" {
		port = "8080"
	}
	logrus.Infof("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}
}

// Helper function to ensure handlers package is imported
// func importHandlers() {}
