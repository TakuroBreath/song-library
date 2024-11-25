package main

import (
	"fmt"
	_ "github.com/TakuroBreath/song-library/docs"
	"github.com/TakuroBreath/song-library/internal/api/handlers"
	"github.com/TakuroBreath/song-library/internal/api/routes"
	"github.com/TakuroBreath/song-library/internal/service"
	"github.com/TakuroBreath/song-library/internal/storage/postgresql"
	"github.com/TakuroBreath/song-library/pkg/migrator"
	"github.com/TakuroBreath/song-library/pkg/sl"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "production"
)

// @title           Song Library API
// @version         1.0
// @description     API Server for Song Library Application
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.basic  BasicAuth
func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	env := os.Getenv("ENV")
	log := setupLogger(env)

	log.Info("starting song-library", slog.String("env", env))
	log.Debug("debug messages are enabled")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	if err := migrator.Migrate(user, password, host, port, dbname, log); err != nil {
		log.Error("failed to apply migrations", sl.Err(err))
		os.Exit(1)
	}

	storage, err := postgresql.NewStorage(psqlInfo, log)
	if err != nil {
		log.Error("failed to create storage", sl.Err(err))
		os.Exit(1)
	}

	songService := service.NewSongService(storage, os.Getenv("API_URL"), log)
	songHandler := handlers.NewSongHandler(songService)

	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	routes.SetupSongRoutes(router, songHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(":8080"); err != nil {
		log.Error("failed to start server", sl.Err(err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
