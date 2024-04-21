package main

import (
	"context"
	"fmt"
	"log"

	"github.com/amitrei/oas-compare-server/database"
	"github.com/amitrei/oas-compare-server/http"
	"github.com/amitrei/oas-compare-server/logger"
)

func main() {

	redis := database.RedisClient{}
	databaseClient := redis.NewClient(&database.ClientConfiguration{
		Address:  "localhost:6379",
		Password: "",
		Context:  context.Background(),
	})

	err := databaseClient.HealthCheck()
	if err != nil {
		errMessage := fmt.Sprintf("Failed initializing the database due to the following error: %s", err.Error())
		log.Fatal(errMessage)
	}

	logger.InitLogger()
	database.InitDatabaseClient(&databaseClient)
	http.NewHttpServer()
}
