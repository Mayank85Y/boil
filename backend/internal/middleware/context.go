package middleware

import (
	"context"

	"github.com/Mayank85Y/boil/internal/logger"
	"github.com/Mayank85Y/boil/internal/server"
	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
)

const (
	UserIDKey   = "user_id"
	UserRoleKey = "user_role"
	LoggerKey   = "Logger"
)

type ContextEnhancer struct {
	server *server.Server
}

func NewContextEnhancer(s *server.Server) *ContextEnhancer{
	return &ContextEnhancer{server: s}
}

func (ce *ContextEnhancer) Enhacer() echo.MiddlewareFunc{
	return func(next echo.HandlerFunc) echo.HandlerFunc{
		return func(c echo.Context) error {
			//extract request id
			requestID := GetRequestID(c)

			contextLogger := ce.server.Logger.With().
				Str("request_id", requestID).
				Str("method", c.Request().Method).
				Str("path", c.Path()).
				Str("ip", c.RealIP()).
				Logger()

			//trace context / transaction 
			if txn := newrelic.FromContext(c.Request().Context()); txn != nil{
				contextLogger = logger.WithTraceContext(contextLogger, txn)
			}

			//extract user info from jwt token
			if userID := ce.extractUserID(c); userID != "" {
				contextLogger = contextLogger.With().Str("user_id", userID).Logger()
			}

			if userRole := ce.extractUserRole(c); userRole != ""{
				contextLogger = contextLogger.With().Str("user_role", userRole).Logger()
			}

			c.Set(LoggerKey, &contextLogger) //store enhanced logger in context

			//createe a new context wih logger
			ctx := context.WithValue(c.Request().Context(), LoggerKey, &contextLogger)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func (ce *ContextEnhancer) extractUserID(c echo.Context) string{
	//check if userid set bu auth middleware
	if userID, ok := c.Get("user_id").(string); ok && userID != ""{
		return userID
	}
	return ""
}

func (ce *ContextEnhancer) extractUserRole(c echo.Context) string{
	//check if userrole set bu auth middleware
	if userRole, ok := c.Get("user_role").(string); ok && userRole != ""{
		return userRole
	}
	return ""
}

func GetUserID(c echo.Context) string{
	if userID, ok := c.Get(UserIDKey).(string); ok{
		return userID
	}
	return ""
}

func GetLogger( c echo.Context) *zerolog.Logger{
	if logger, ok := c.Get(LoggerKey).(*zerolog.Logger); ok {
		return logger
	}
	logger := zerolog.Nop()
	return &logger
}