package rabbitmq

import (
	"context"
	"encoding/json"
	
	amqp "github.com/rabbitmq/amqp091-go"
)

type Headers amqp.Table

// ProduceWithContext message should be json serialized
func ProduceWithContext(ctx context.Context, ch *amqp.Channel, message any, headers Headers, exchange, key string) error {
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
