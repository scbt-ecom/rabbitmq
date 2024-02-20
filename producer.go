package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"reflect"
)

// ProduceWithContext message should be structure pointer
func ProduceWithContext(ctx context.Context, ch *amqp.Channel, message any, headers amqp.Table, exchange, key string) error {
	v := reflect.ValueOf(message)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("message not a structure pointer")
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx, exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		Headers:     headers,
	})
	if err != nil {
		return err
	}

	return nil
}
