package job

import (
	"context"
	"log"
	"math/rand"
	"time"
)

type DefaultHandler struct{}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}

func (h *DefaultHandler) Handle(ctx context.Context, job Job) error {
	log.Printf(
		"Executing job %s (type=%s, payload=%s)",
		job.ID,
		job.Type,
		job.Payload,
	)

	// job processing logic here
	// Simulate job processing time with a randon value between 1 and 3 seconds
	time.Sleep(time.Duration(1+rand.Intn(3)) * time.Second)

	return nil
}
