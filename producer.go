package rabbitmq

import (
	"encoding/json"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"reflect"
)

// Produce message should be structure pointer
func Produce(conn *amqp.Connection, message any, headers amqp.Table, exchange, key string) error {
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

	err = ch.Publish(exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		Headers:     headers,
	})
	if err != nil {
		return err
	}

	return nil
}
