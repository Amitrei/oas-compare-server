package middlewares

import (
	"github.com/amitrei/oas-compare-server/logger"
	"github.com/labstack/echo"
)

func HttpLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			res := c.Response()
			log := logger.GetContextLogger(c)

			n := res.Status
			switch {
			case n >= 500:
				log.Error("Server error")
			case n >= 400:
				log.Warn("Client error")
			case n >= 300:
				log.Info("Redirection")
			default:
				log.Info("Success")
			}

			return nil
		}
	}
}
