package main

import (
	"context"
	"fmt"
	"time"

	"github.com/scbt-ecom/rabbitmq"
)

func main() {
	publisher := rabbitmq.NewPublisher("amqp://guest:guest@localhost:5672/", "testchannel", "dispatcher", "controller")

	go publisher.StartLifeSupport(10 * time.Second)

	for i := 0; i < 10; i++ {
		go func() {
			for {
				time.Sleep(10 * time.Millisecond)
				err := rabbitmq.ProduceWithContext(context.Background(), publisher.GetChannel("testchannel"), map[string]interface{}{"sad": "l;ol"}, rabbitmq.Headers{}, "amq.direct", "test")
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}()
	}

	select {}
}
