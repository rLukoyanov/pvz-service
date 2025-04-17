package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware возвращает Echo middleware с логированием через logrus
func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			res := c.Response()

			err := next(c)

			stop := time.Now()
			latency := stop.Sub(start)

			entry := logrus.WithFields(logrus.Fields{
				"time":       stop.Format(time.RFC3339),
				"remote_ip":  c.RealIP(),
				"host":       req.Host,
				"method":     req.Method,
				"uri":        req.RequestURI,
				"user_agent": req.UserAgent(),
				"status":     res.Status,
				"latency":    latency,
				"latency_ms": float64(latency.Microseconds()) / 1000.0,
				"request_id": res.Header().Get(echo.HeaderXRequestID),
			})

			if err != nil {
				entry = entry.WithField("error", err)
				entry.Error("request failed")
			} else {
				entry.Info("request handled")
			}

			return err
		}
	}
}
