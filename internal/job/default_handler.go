package job

import (
	"context"
	"log"
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

	// sleep for 5 seconds to simulate job processing
	time.Sleep(5 * time.Second)

	return nil
}
