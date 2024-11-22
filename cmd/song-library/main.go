package main

import (
	"fmt"
	"github.com/TakuroBreath/song-library/internal/storage/postgresql"
	"github.com/TakuroBreath/song-library/pkg/migrator"
	"github.com/joho/godotenv"
	"os"
)

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

	_ = storage
}
