package handler

import (
	"github.com/Mayank85Y/boil/internal/server"
	"github.com/Mayank85Y/boil/internal/service"
)

type Handlers struct{
	Health 	*HealthHandler
	OpenAPI	*OpenAPIHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers{
	return &Handlers{
		Health:  NewHealthHandler(s),
		OpenAPI: NewOpenAPIHandler(s),
	}
}
