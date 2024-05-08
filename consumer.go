package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func (q *Queue) Consume(conn *amqp.Connection, consumerName string) (<-chan amqp.Delivery, error) {
	ch, err := conn.Channel()
	if err != nil {
		ch.Close()
		return nil, err
	}

	msgs, err := ch.Consume(q.Name, consumerName, false, false, false, false, nil)
	if err != nil {
		ch.Close()
		return nil, err
	}

	return msgs, nil
}
