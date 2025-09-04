package repository

import "github.com/Mayank85Y/boil/internal/server"

type Repositories struct{}

func NewRepositories(s *server.Server) *Repositories{
	return &Repositories{}
}