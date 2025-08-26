package service

import (
	"github.com/Mayank85Y/boil/internal/lib/job"
	"github.com/Mayank85Y/boil/internal/repository"
	"github.com/Mayank85Y/boil/internal/server"
)

type Services struct {
	Auth *AuthService
	Job  *job.JobService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error){
	authService := NewAuthService(s)

	return &Services{
		Auth: 	authService,
		Job: 	s.Job,
	}, nil
}