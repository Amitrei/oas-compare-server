package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/amitrei/oas-compare-server/database"
	"github.com/amitrei/oas-compare-server/handlers"
	"github.com/amitrei/oas-compare-server/logger"
	"github.com/labstack/echo"
)

func NewHttpServer() {
	e := echo.New()
	e.HTTPErrorHandler = handlers.GlobalErrorHandler

	group := e.Group("/oas-report-server/v1")
	group.Use(logger.HttpLogger())
	for _, h := range handlers.GetHandlers() {
		group.Add(h.Method, h.Path, h.HandlerFunc)
	}

	group.GET("/health", healthCheck)
	serverPort := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	e.Logger.Fatal(e.Start(serverPort))
}

func healthCheck(ctx echo.Context) error {
	err := database.GetDatabaseClient().HealthCheck()
	if err != nil {
		logger := logger.GetContextLogger(ctx)
		errMessage := fmt.Sprintf("Failed initializing the database due to the following error: %s", err.Error())
		logger.Error(errMessage)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusOK)

}
