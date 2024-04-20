package http

import (
	"github.com/amitrei/oas-compare-server/handlers"
	"github.com/amitrei/oas-compare-server/http/middlewares"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

var logger, _ = zap.NewProduction()

func NewHttpServer() {
	e := echo.New()
	e.Use(middlewares.HttpLogger(logger))
	for _, h := range handlers.GetHandlers() {
		e.Add(h.Method, h.Path, h.HandlerFunc)
	}
	e.Logger.Fatal(e.Start(":8080"))
}
