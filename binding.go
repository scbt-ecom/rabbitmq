package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func initBinding(conn *amqp.Connection, queue Queue) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.QueueBind(queue.Name, queue.Key, queue.Exchange, false, nil)
	if err != nil {
		return err
	}

	return nil
}
