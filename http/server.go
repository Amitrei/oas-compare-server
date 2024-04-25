package http

import (
	"fmt"
	"os"

	"github.com/amitrei/oas-compare-server/handlers"
	"github.com/amitrei/oas-compare-server/logger"
	"github.com/labstack/echo"
)

func NewHttpServer() {
	e := echo.New()
	e.Use(logger.HttpLogger())
	e.HTTPErrorHandler = handlers.GlobalErrorHandler

	for _, h := range handlers.GetHandlers() {
		e.Add(h.Method, h.Path, h.HandlerFunc)
	}

	serverPort := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	e.Logger.Fatal(e.Start(serverPort))

}
