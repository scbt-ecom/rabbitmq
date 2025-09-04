package rabbitmq

import "testing"

func TestInitQueue(t *testing.T) {
	q := Queue{
		Name:     "test",
		Key:      "test",
		Exchange: "amq.direct",
	}

	err := InitQueue(conn, q)
	if err != nil {
		t.Fatalf("failed to init queue: %s", err.Error())
	}
}
