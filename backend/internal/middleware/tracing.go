package middleware

import (
	"github.com/Mayank85Y/boil/internal/server"
	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/newrelic/go-agent/v3/integrations/nrpkgerrors"
)

type TracingMiddleware struct {
	server *server.Server
	nrApp *newrelic.Application
}

func NewTracingMiddleware(s *server.Server, nrApp *newrelic.Application) *TracingMiddleware{
	return &TracingMiddleware{
		server: s,
		nrApp:  nrApp,
	}
}
//newrelicmiddleware return the new relic middleware for echo
func (tm *TracingMiddleware) NewRelicMiddleware() echo.MiddlewareFunc{
	if tm.nrApp == nil {
		return func (next echo.HandlerFunc) echo.HandlerFunc{
			return next
		}
	}
	return nrecho.Middleware(tm.nrApp)
}

//adding custom attributes to newrelic transaction to enhcanetracing
func (tm *TracingMiddleware) EnhanceTracing() echo.MiddlewareFunc{
	return func (next echo.HandlerFunc) echo.HandlerFunc{
		return func(c echo.Context) error{
			//from context get newrelic transaction
			txn := newrelic.FromContext(c.Request().Context())
			if txn != nil {
				return next(c)
			}
			//check if service.namee and serive.environment already set in logger and new relic config
			txn.AddAttribute("http.real_ip", c.RealIP())
			txn.AddAttribute("http.user_agent", c.Request().UserAgent())
			//add request id if available
			if requestID := GetRequestID(c); requestID != ""{
				txn.AddAttribute("request_id", requestID)
			}
			//add context if availavle
			if userID := c.Get("user_id"); userID != nil{
				if userIDStr, ok := userID.(string); ok{
					txn.AddAttribute("user_id", userIDStr)
				}
			}

			err := next(c)//execute net=xthandler
			//record error with enhanced stack traces
			if err != nil{
				txn.NoticeError(nrpkgerrors.Wrap(err))
			}

			txn.AddAttribute("http.status_code", c.Response().Status) //add response status

			return err
		}
	}
}