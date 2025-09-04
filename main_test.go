package rabbitmq

import (
	"os"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection

func TestMain(m *testing.M) {
	newConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	conn = newConn

	os.Exit(m.Run())
}
