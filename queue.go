package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	Name                string
	Key                 string
	Exchange            string
	XDeadLetterExchange string
	XMessageTTL         int //milliseconds
}

func InitQueue(conn *amqp.Connection, queue Queue) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	args := amqp.Table{}
	if queue.XDeadLetterExchange != "" {
		args["x-dead-letter-exchange"] = queue.XDeadLetterExchange
	}
	if queue.XMessageTTL != 0 {
		args["x-message-ttl"] = queue.XMessageTTL
	}

	_, err = ch.QueueDeclare(queue.Name, true, false, false, false, args)
	if err != nil {
		return err
	}

	err = initBinding(conn, queue)
	if err != nil {
		return err
	}

	return nil
}
