package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"reflect"
    "sync"
    "fmt"
)

type Headers amqp.Table

func ProduceWithContext(ctx context.Context, ch *amqp.Channel, message any, headers Headers, exchange, key string) error {
	v := reflect.ValueOf(message)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("message must be a pointer to a structure")
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Включаем подтверждения доставки.
	err = ch.Confirm(false)
	if err != nil {
		return fmt.Errorf("failed to set channel to confirm mode: %w", err)
	}

	ackCh := ch.NotifyPublish(make(chan amqp.Confirmation, 100))
	returnCh := ch.NotifyReturn(make(chan amqp.Return, 100))

	var wg sync.WaitGroup
	wg.Add(1)

	var publishErr error
	go func() {
		defer wg.Done()
		select {
		case ret := <-returnCh:
			publishErr = fmt.Errorf("failed to route message: %s", string(ret.Body))
		case confirm := <-ackCh:
			if !confirm.Ack {
				publishErr = fmt.Errorf("message not acknowledged by broker")
			}
		case <-ctx.Done():
			publishErr = fmt.Errorf("context deadline exceeded: %w", ctx.Err())
		}
	}()

	err = ch.PublishWithContext(ctx, exchange, key, true, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		Headers:     amqp.Table(headers),
	})
	if err != nil {
		return fmt.Errorf("publish error: %w", err)
	}

	wg.Wait()
	return publishErr
}
