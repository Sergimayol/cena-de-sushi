package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	TO_PRODUCE = 10
)

var (
	SHUSHI_TYPES = []string{"nigiris de salmó", "sashimis de tonyina", "makis de cranc"}
)

type Empty struct{}

type Shushi struct {
	ShushiType string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func producir(ch *amqp.Channel, q amqp.Queue) {
	// Sacar del canal que el canal esta vacio
	for i := 0; i < TO_PRODUCE; i++ {
		body := SHUSHI_TYPES[i%len(SHUSHI_TYPES)]
		err := ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s\n", body)
	}
	// Avisar al canal de que el canal esta lleno
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	producir(ch, q)
}
