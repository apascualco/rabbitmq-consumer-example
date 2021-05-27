package main

import (
	"bytes"
	"log"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

const QUEUE = "example_task_queue"

func connectRabbitMQ() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/apascualco")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	return conn
}

func openChannel(conn *amqp.Connection) *amqp.Channel {
	c, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	return c
}

func declareQueue(channel *amqp.Channel) amqp.Queue {
	q, err := channel.QueueDeclare(QUEUE, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}
	return q
}

func publish(text string, channel *amqp.Channel, queue amqp.Queue) {
	err := channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(text),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}
}

func generateMessages(count int, channel *amqp.Channel, queue amqp.Queue) {
	n := 0
	for n < count {
		n++
		msg := "apascualco count (" + strconv.Itoa(n) + "): " + time.Now().String()
		log.Printf("Publishing: %s", msg)
		publish(msg, channel, queue)
		time.Sleep(2 * time.Second)
	}
}

func qualityOfService(channel *amqp.Channel) {
	err := channel.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("Failed to configure channel with quality of services: %s", err)
	}
}

func consumer(channel *amqp.Channel, queue amqp.Queue) <-chan amqp.Delivery {
	consumer, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to configure consumer")
	}
	return consumer
}

func main() {
	conn := connectRabbitMQ()
	defer conn.Close()

	channel := openChannel(conn)
	defer channel.Close()

	queue := declareQueue(channel)

	log.Println("Generating test messages")
	generateMessages(10, channel, queue)

	qualityOfService(channel)

	log.Println("Consuming all messages")
	chanDelivery := consumer(channel, queue)
	for d := range chanDelivery {
		log.Printf("Consuming: %s", d.Body)
		dotCount := bytes.Count(d.Body, []byte("."))
		t := time.Duration(dotCount)
		time.Sleep(t * time.Second)
		d.Ack(false)
	}
}
