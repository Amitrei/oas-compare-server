package database

import (
	"context"
)

var databaseClientRef DatabaseClient

func InitDatabaseClient(databaseClient *DatabaseClient) {
	if databaseClientRef == nil {
		databaseClientRef = *databaseClient
	}
}

func GetDatabaseClient() DatabaseClient {
	return databaseClientRef
}

type DatabaseClient interface {
	NewClient(clientConfig *ClientConfiguration) DatabaseClient
	HealthCheck() error
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
}

type ClientConfiguration struct {
	Address  string
	Password string
	Context  context.Context
}
