package job

import (
	"context"
	"log"
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

	return nil
}
