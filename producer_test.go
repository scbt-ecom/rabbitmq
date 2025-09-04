package rabbitmq

import (
	"context"
	"testing"
)

func TestProduceWithContext(t *testing.T) {
	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("failed to open channel: %s", err.Error())
	}

	msg := map[string]interface{}{
		"product_name": "КК Халва",
	}

	for i := 0; i < needMessagesForSuccessful; i++ {
		err = ProduceWithContext(context.Background(), ch, msg, Headers{}, "amq.direct", "test")
		if err != nil {
			t.Fatalf("failed to produce message: %s", err.Error())
		}
	}
}
