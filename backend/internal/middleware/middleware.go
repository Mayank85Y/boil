package middleware

import (
	"github.com/Mayank85Y/boil/internal/server"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Middlewares struct {
	Global          *GlobalMiddlewares
	Auth            *AuthMiddleware
	ContextEnhancer *ContextEnhancer
	Tracing         *TracingMiddleware
	RateLimit       *RateLimitMiddleware
}

func NewMiddlewares(s *server.Server) *Middlewares{
	var nrApp *newrelic.Application  //newrelic intance from server
	if s.LoggerService != nil {
		nrApp = s.LoggerService.GetApplication()
	}
		//this pattern is known as dependecy injection	
		return &Middlewares{
			Global:				NewGlobalMiddlewares(s),
			Auth:				NewAuthMiddleware(s),			
			ContextEnhancer: 	NewContextEnhancer(s),
			Tracing:			NewTracingMiddleware(s, nrApp),
			RateLimit:			NewRateLimitMiddleware(s),
		}
}