package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
	"github.com/skbt-ecom/logging"
)

type Exchange struct {
	Name string
	Type string
}

func InitExchange(conn *amqp091.Connection, log *logging.Logger, exchange Exchange) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	err = ch.ExchangeDeclare(exchange.Name, exchange.Type, true, false, false, false, nil)
	if err != nil {
		return err
	}

	return nil
}
