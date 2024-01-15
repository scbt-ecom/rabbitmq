package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateConnection(rabbitmqUrl string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(rabbitmqUrl)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
