package main

import (
	"context"
	"fmt"
	"log"

	"github.com/amitrei/oas-compare-server/database"
	"github.com/amitrei/oas-compare-server/http"
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
	database.DatabaseClientRef = databaseClient

	http.NewHttpServer()
}
