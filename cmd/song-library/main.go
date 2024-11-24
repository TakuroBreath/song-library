package main

import (
	"fmt"
	_ "github.com/TakuroBreath/song-library/docs"
	"github.com/TakuroBreath/song-library/internal/api/handlers"
	"github.com/TakuroBreath/song-library/internal/api/routes"
	"github.com/TakuroBreath/song-library/internal/service"
	"github.com/TakuroBreath/song-library/internal/storage/postgresql"
	"github.com/TakuroBreath/song-library/pkg/migrator"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
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

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	migrator.Migrate(user, password, host, port, dbname)

	storage, err := postgresql.NewStorage(psqlInfo)
	if err != nil {
		panic(err)
	}

	songService := service.NewSongService(storage, os.Getenv("API_URL"))
	songHandler := handlers.NewSongHandler(songService)
	router := gin.Default()

	routes.SetupSongRoutes(router, songHandler)

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
