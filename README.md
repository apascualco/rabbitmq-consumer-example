# rabbitmq-consumer-example

## How to test

Launch docker with RabbitMQ

```
docker-composer up -d
```

After that build the example and run

```
go build cmd/consumer_example.go

./consumer_example
```

To exit after result you should use "ctrl+c"

Remember rabbitmq admin url http://localhost:15672/ (guest:guest)