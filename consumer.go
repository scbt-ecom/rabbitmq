package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func Consume(conn *amqp.Connection, log *logrus.Entry, queueName, consumerName string) (<-chan amqp.Delivery, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(queueName, consumerName, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println(msgs)
	return msgs, nil
}
