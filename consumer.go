package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func Consume(conn *amqp.Connection, queueName, consumerName string) (<-chan amqp.Delivery, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	msgs, err := ch.Consume(queueName, consumerName, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
