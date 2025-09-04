package job

import (	
	"github.com/hibiken/asynq"
	"encoding/json"
	"time"
	)

const (
	TaskWelcome = "email: welcome"
)

type WelcomeEmailPayLoad struct {
	To        string `json:"to"`
	FirstName string `json:"first_name"`
}

func NewWelcomeEmailTask(to string, firstName string) (*asynq.Task, error) {
	payload, err := json.Marshal(WelcomeEmailPayLoad{
		To:		   to,
		FirstName: firstName,
	})
	if err != nil{
		return nil, err
	}

	return asynq.NewTask(TaskWelcome, payload,
		asynq.MaxRetry(3),
		asynq.Queue("default"),
		asynq.Timeout(30*time.Second)), nil
}