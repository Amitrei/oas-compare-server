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
