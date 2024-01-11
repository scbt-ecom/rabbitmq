package rabbitmq

import (
	"encoding/json"
	"errors"
	"github.com/rabbitmq/amqp091-go"
	"reflect"
)

// Produce message should be structure pointer
func Produce(conn *amqp091.Connection, message any, exchange, key string) error {
	//fmt.Println(message)
	v := reflect.ValueOf(message)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("message not a structure pointer")
	}
	//fmt.Println(message)
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	//fmt.Println(message)
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	//fmt.Println(body)
	err = ch.Publish(exchange, key, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return err
	}

	return nil
}
