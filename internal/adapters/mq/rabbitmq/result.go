package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"speed-typing-auth-service/pkg/logging"
)

const queueName = "results"

type ResultsMigrator struct {
	channel *amqp.Channel
	log     logging.Logger
}

func NewResultsMigrator(channel *amqp.Channel, log logging.Logger) *ResultsMigrator {
	_, err := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatal("failed to declare queue")
	}
	return &ResultsMigrator{
		channel: channel,
		log:     log,
	}
}

func (r *ResultsMigrator) MigrateSessionResults(session string, userID string) {
	if session != "" {
		message, err := json.Marshal(
			map[string]interface{}{
				"session": session,
				"userID":  userID,
			},
		)
		if err != nil {
			r.log.Warnf(`results not migrated due to error: %v`, err)
			return
		}
		err = r.channel.Publish(
			"",
			queueName,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        message,
			},
		)
		if err != nil {
			r.log.Warnf(`results not migrated due to error: %v`, err)
			return
		}
	}
}
