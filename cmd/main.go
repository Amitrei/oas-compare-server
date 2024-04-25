package main

import (
	"context"
	"fmt"
	"os"

	"github.com/amitrei/oas-compare-server/database"
	"github.com/amitrei/oas-compare-server/http"
	"github.com/amitrei/oas-compare-server/logger"
	"github.com/joho/godotenv"
)

func main() {

	logger.InitLogger()
	logger := logger.GetLogger()
	err := godotenv.Load()
	if err != nil {
		warningMessage := fmt.Sprintf("Failed loading .env file due to the following error: %s", err)
		logger.Warn(warningMessage)
	}

	redis := database.RedisClient{}
	databaseClient := redis.NewClient(&database.ClientConfiguration{
		Address:  os.Getenv("DATABASE_ADDRESS"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Context:  context.Background(),
	})

	err = databaseClient.HealthCheck()
	if err != nil {
		errMessage := fmt.Sprintf("Failed initializing the database due to the following error: %s", err.Error())
		logger.Fatal(errMessage)
	}

	database.InitDatabaseClient(&databaseClient)
	http.NewHttpServer()
}
