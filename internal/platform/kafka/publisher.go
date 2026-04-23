package kafka

import (
	"context"
	"log"

	"github.com/biangacila/kopesa-loan-platform-api/internal/shared"
)

type Publisher interface {
	Publish(ctx context.Context, topic string, event shared.Event) error
}

type NoopPublisher struct{}

func NewNoopPublisher() *NoopPublisher {
	return &NoopPublisher{}
}

func (p *NoopPublisher) Publish(_ context.Context, topic string, event shared.Event) error {
	log.Printf("noop kafka publish topic=%s type=%s id=%s", topic, event.Type, event.ID)
	return nil
}
