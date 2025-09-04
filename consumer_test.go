package rabbitmq

import (
	"context"
	"testing"
	"time"
)

const needMessagesForSuccessful = 5

func TestConsume(t *testing.T) {
	q := Queue{
		Name:     "test",
		Key:      "test",
		Exchange: "amq.direct",
	}

	ch, err := conn.Channel()
	if err != nil {
		t.Fatal(err)
	}

	err = InitQueue(conn, q)
	if err != nil {
		t.Fatalf("init queue failed: %s", err.Error())
	}

	msgs, err := q.Consume(ch, "test-consumer")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var counter int
	for counter < needMessagesForSuccessful {
		select {
		case <-ctx.Done():
			t.Fatalf("timeout: received %d, expected %d", counter, needMessagesForSuccessful)
		case msg, ok := <-msgs:
			if !ok {
				t.Fatalf("channel is closed: received %d, expected %d", counter, needMessagesForSuccessful)
			}
			counter++
			msg.Ack(false)
		}

	}

}
