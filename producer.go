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
func ProduceWithContext(ctx context.Context, ch *amqp.Channel, message string, exchange, key string) error {
    var wg sync.WaitGroup
    wg.Add(1)
    returnCh := make(chan amqp.Return, 1)
    ch.NotifyReturn(returnCh)
    ch.Confirm(false)
    ackCh := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
    var publishErr error
    var receivedReturn bool
    go func() {
        defer wg.Done()
        select {
        case ret := <-returnCh:
            receivedReturn = true
            publishErr = fmt.Errorf("failed to route message: %s", string(ret.Body))
        case confirm := <-ackCh:
            if !confirm.Ack {
                publishErr = fmt.Errorf("message not acknowledged by broker")
            }
        case <-ctx.Done():
            if !receivedReturn {
                publishErr = fmt.Errorf("context deadline exceeded, but no return received")
            }
        }
    }()
    err := ch.PublishWithContext(ctx, exchange, key, true, false, amqp.Publishing{
        ContentType: "text/plain",
        Body:        []byte(message),
    })
    if err != nil {
        return fmt.Errorf("publish error: %w", err)
    }
    wg.Wait()
    return publishErr
}
