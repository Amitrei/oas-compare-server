package http

import (
	"github.com/amitrei/oas-compare-server/handlers"
	"github.com/amitrei/oas-compare-server/http/middlewares"
	"github.com/labstack/echo"
)

func NewHttpServer() {
	e := echo.New()

	e.Use(middlewares.HttpLogger())
	for _, h := range handlers.GetHandlers() {
		e.Add(h.Method, h.Path, h.HandlerFunc)
	}

	e.Logger.Fatal(e.Start(":8080"))

}
