# RabbitMQ package for Совкомбанк Технологии

## Getting started
```bash
go get github.com/skbt-ecom/rabbitmq@v1.3.3
```

## Development

### Create connection with RabbitMQ || Enable life supporting on connection/channels
```
pub := rabbitmq.NewPublisher("amqp://guest:guest@localhost:5672/", "dispatcher", "controller")

go pub.StartLifeSupport(10 * time.Second)

for {
	time.Sleep(10 * time.Millisecond)
	err := rabbitmq.ProduceWithContext(context.Background(), pub.GetChannel("controller"), map[string]interface{}{"sad": "l;ol"}, rabbitmq.Headers{}, "amq.direct", "test")
	if err != nil {
		fmt.Println(err.Error())
	}
}
```
### Example of queue initialization

````
exampleQueue := rabbitmq.Queue{
		Name:                cfg.RabbitMqQueueName, // required
		Key:                 cfg.RabbitMqKey, // required
		Exchange:            cfg.RabbitMqExchange, // required
		XDeadLetterExchange: "retry", // optionally
		XMessageTTL:         1000 * 60 * 60, // optionally
	}

err := rabbitmq.InitQueue(conn, exampleQueue)
````

### Example of exchange initialization
````
exampleExchange := rabbitmq.Exchange{
        Name: cfg.ExchangeName,
        Type: cfg.ExchangeType,
    }
    
err := rabbitmq.InitExchanges(conn, exampleExchange)
````

### Example of queue consuming
````
msgs, err := rabbitmq.Consume(conn, <queueName>, <consumerName>)
````

### Example of queue producing
````
err := rabbitmq.ProduceWithContext(ctx, pub.GetChannel("dispatcher"), <messageStructure>, <headers>, <exchangeName>, <key>)
// messageStructure should be a structure pointer
````