package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func initBinding(conn *amqp.Connection, log *logrus.Entry, queue Queue) error {
	ch, err := conn.Channel()
	if err != nil {
		log.WithField("error", err).Error("failed to open channel")
		return err
	}

	err = ch.QueueBind(queue.Name, queue.Name, queue.Exchange, false, nil)
	if err != nil {
		log.WithField("error", err).Error("failed to bind queue")
		return err
	}

	return nil
}
