# RabbitMQ package for Совкомбанк Технологии

## Getting started
```bash
go get github.com/skbt-ecom/rabbitmq
```

## Development

### Create connection with RabbitMQ
```
conn, err := rabbitmq.CreateConnection(<rabbitMqUrl>)
// url example - amqp://<username>:<password>@<host>:<port>
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
err := rabbitmq.ProduceWithContext(ctx, ch, <messageStructure>, <headers>, <exchangeName>, <key>)
// messageStructure should be a structure pointer
````