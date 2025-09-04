package middleware

import "github.com/Mayank85Y/boil/internal/server"

type RateLimitMiddleware struct {
	server *server.Server
}

func NewRateLimitMiddleware(s *server.Server) *RateLimitMiddleware{
	return &RateLimitMiddleware{
		server: s,
	}
}

func (r *RateLimitMiddleware) RecordRateLimitHit(endpoint string){
	if r.server.LoggerService != nil && r.server.LoggerService.GetApplication() != nil{
		r.server.LoggerService.GetApplication().RecordCustomEvent("RateLimitHi", map[string]interface{}{
			"endpoint": endpoint,
		})
	}
}