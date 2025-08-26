package server

import (
	"fmt"
	"net/http"
	"errors"
	"time"
	"context"

	"github.com/Mayank85Y/boil/internal/config"
	"github.com/Mayank85Y/boil/internal/database"
	"github.com/Mayank85Y/boil/internal/lib/job"
	loggerPkg "github.com/Mayank85Y/boil/internal/logger"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/newrelic/go-agent/v3/integrations/nrredis-v9"
)

type Server struct {
	Config			*config.Config
	Logger			*zerolog.Logger 
	LoggerService 	*loggerPkg.LoggerService
	DB              *database.Database
	Redis 			*redis.Client
	httpServer      *http.Server
	Job 			*job.JobService
}

func New(cfg *config.Config, logger *zerolog.Logger, loggerService *loggerPkg.LoggerService) (*Server, error){
	db, err := database.New(cfg, logger, loggerService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	//redisclient with new relic integration
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
	})

	//add new relic redis hook (if available)
	if loggerService != nil && loggerService.GetApplication() != nil{
		redisClient.AddHook(nrredis.NewHook(redisClient.Options()))
	}

	//test redis conn
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil{
		logger.Error().Err(err).Msg("Failed to connect to redis, continuing without redis")
	}

	//job service
	jobService := job.NewJobService(logger, cfg)
	jobService.InitHandlers(cfg, logger)

	//start job server
	if err := jobService.Start(); err != nil{
		return nil, err
	}

	server := &Server{
		Config: 		cfg,
		Logger: 		logger,
		LoggerService: 	loggerService,
		DB:             db,
		Redis: 			redisClient,	
		Job: 			jobService,
	}

	//runtime metrics are auto collected by newrelic
	return server, nil
}

func (s *Server) SetupHTTPServer(handler http.Handler){
	s.httpServer = &http.Server{
		Addr:         ":" + s.Config.Server.Port,
		Handler:	  handler,
		ReadTimeout:  time.Duration(s.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.Config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.Config.Server.IdleTimeout) * time.Second,
	}
}

func (s *Server) Start() error{
	if s.httpServer == nil {
		return errors.New("HTTP server not initialized")
	}

	s.Logger.Info().
		Str("port", s.Config.Server.Port).
		Str("env", s.Config.Primary.Env).
		Msg("Starting server")
		
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error{
	if err := s.httpServer.Shutdown(ctx); err != nil{
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	if err := s.DB.Close(); err != nil{
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	if s.Job != nil{
		s.Job.Stop()
	}
	return nil
}
