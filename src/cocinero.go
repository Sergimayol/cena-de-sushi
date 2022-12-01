package main

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	TO_PRODUCE    = 10
	RABBITMQ_HOST = "amqp://guest:guest@localhost:5672/"
)

var (
	SHUSHI_TYPES = []string{"nigiris de salmó", "sashimis de tonyina", "makis de cranc"}
	shushis      [TO_PRODUCE]string
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func producir() {
	for i := 0; i < TO_PRODUCE; i++ {
		shushi := SHUSHI_TYPES[i%len(SHUSHI_TYPES)]
		shushis[i] = fmt.Sprintf("Produciendo shushi: %s", shushi)
		log.Printf(" [x] Sent %s\n", shushi)
	}
}

func enviar(ch *amqp.Channel, ctx context.Context, q amqp.Queue, err error) {
	for i := 0; i < len(shushis); i++ {
		err = ch.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(shushis[i]),
			})
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s", shushis[i])
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial(RABBITMQ_HOST)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	producir()

	enviar(ch, ctx, q, err)

}
