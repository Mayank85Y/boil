package job

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Mayank85Y/boil/internal/config"
	"github.com/Mayank85Y/boil/internal/lib/email"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

var emailClient *email.Client

func (j *JobService) InitHandlers(config *config.Config, logger *zerolog.Logger){
	emailClient = email.NewClient(config, logger)
}

func (j *JobService) handleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var p WelcomeEmailPayLoad
	if err := json.Unmarshal(t.Payload(), &p); err != nil{
		return fmt.Errorf("failed to unmarshal welcome email payLoad: %w", err)
	}

	j.logger.Info().
		Str("type", "welcome").
		Str("to", p.To).
		Msg("Sending welcome email")

	err := emailClient.SendWelcomeEmail(
		p.To,
		p.FirstName,
	)
	if err != nil{
		j.logger.Error().
			Str("type", "welcome").
			Str("to", p.To).
			Err(err).
			Msg("Failed to send welcome email")
		return err
	}

	j.logger.Info().
		Str("type", "welcome").
		Str("to", p.To).
		Msg("Successfully sent welcome email")

	return nil
}