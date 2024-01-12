package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/skbt-ecom/logging"
)

func CreateConnection(rabbitmqUrl string, log *logging.Logger) (*amqp.Connection, error) {
	conn, err := amqp.Dial(rabbitmqUrl)
	if err != nil {
		log.WithField("error", err).Error("failed to establish connection with rabbitmq")
		return nil, err
	}

	return conn, nil
}
