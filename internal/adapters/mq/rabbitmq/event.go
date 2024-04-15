package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ttodoshi/code-typing-auth-service/internal/core/domain"
	"github.com/ttodoshi/code-typing-auth-service/internal/core/ports"
	"github.com/ttodoshi/code-typing-auth-service/pkg/logging"
	"time"
)

type EventDispatcher struct {
	channel *amqp.Channel
	log     logging.Logger
}

func NewEventDispatcher(channel *amqp.Channel, log logging.Logger) ports.EventDispatcher {
	for _, exchange := range ports.Exchanges {
		err := channel.ExchangeDeclare(
			exchange,
			"fanout",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("failed to declare exchange due to: %s", err.Error())
		}
	}
	return &EventDispatcher{
		channel: channel,
		log:     log,
	}
}

func (r *EventDispatcher) Dispatch(event domain.Event) {
	ctx, cancel := context.WithTimeout(
		context.Background(), 5*time.Second,
	)
	defer cancel()

	err := r.channel.PublishWithContext(
		ctx,
		event.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        event.Body,
		},
	)
	if err != nil {
		r.log.Warnf(`event not dispatched due to error: %v`, err)
		return
	}
}
