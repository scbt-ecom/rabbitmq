package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"reflect"
)

type Headers amqp.Table

// ProduceWithContext message should be structure pointer
func ProduceWithContext(ctx context.Context, ch *amqp.Channel, message any, headers Headers, exchange, key string) error {
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
		Headers:     amqp.Table(headers),
	})
	if err != nil {
		return err
	}

	return nil
}
