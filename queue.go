package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Queue struct {
	Name                string
	Exchange            string
	XDeadLetterExchange string
	XMessageTTL         int //milliseconds
}

func InitQueue(conn *amqp.Connection, log *logrus.Entry, queue Queue) error {
	ch, err := conn.Channel()
	if err != nil {
		log.WithField("error", err).Error("failed to open channel")
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
		log.WithField("error", err).Error("failed to declare queue")
		return err
	}

	err = initBinding(conn, log, queue)
	if err != nil {
		return err
	}

	return nil
}
