package database

import "context"

var DatabaseClientRef DatabaseClient

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
