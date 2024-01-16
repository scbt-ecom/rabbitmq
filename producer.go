package rabbitmq

import (
	"encoding/json"
	"errors"
	"github.com/rabbitmq/amqp091-go"
	"reflect"
)

// Produce message should be structure pointer
func Produce(conn *amqp091.Connection, message any, exchange, key string) error {
	v := reflect.ValueOf(message)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("message not a structure pointer")
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = ch.Publish(exchange, key, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return err
	}

	return nil
}
