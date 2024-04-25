package logger

import (
	"time"

	"github.com/labstack/echo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger() {

	if logger == nil {
		logger, _ = zap.NewProduction()
	}

}

func GetLogger() *zap.Logger {
	return logger
}

func GetContextLogger(ctx echo.Context) *zap.Logger {

	start := time.Now()

	req := ctx.Request()
	res := ctx.Response()

	traceId := req.Header.Get("X-B3-TraceId")
	if traceId == "" {
		traceId = res.Header().Get(echo.HeaderXRequestID)
	}
	fields := []zapcore.Field{
		zap.Int("status", res.Status),
		zap.String("traceId", traceId),
		zap.String("latency", time.Since(start).String()),
		zap.String("method", req.Method),
		zap.String("uri", req.RequestURI),
		zap.String("host", req.Host),
	}

	return logger.With(fields...)
}

func HttpLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			res := c.Response()
			log := GetContextLogger(c)

			n := res.Status
			switch {
			case n >= 500:
				log.Error("Server error")
			case n >= 400:
				log.Error("Client error")
			case n >= 300:
				log.Info("Redirection")
			default:
				log.Info("Success")
			}

			return nil
		}
	}
}
