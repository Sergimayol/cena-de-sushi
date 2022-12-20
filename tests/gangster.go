package main

import (
	"log"
	"math/rand"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	RABBITMQ_HOST = "amqp://guest:guest@localhost:5672/"
)

type Empty struct{}

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

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	aviso, err := ch.QueueDeclare(
		"avisos", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgsAviso, err := ch.Consume(
		aviso.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	gangster, err := ch.QueueDeclare(
		"gangster", // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// The number of shuhis to eat is the length of the queue
	numshuhisAComer := gangster.Messages

	log.Printf("Quiero comer %d piezas de shushi", numshuhisAComer)

	done := make(chan Empty, 1)

	// Receive messages from the queue
	go func() {
		// Get the messages from the queue
		for {
			log.Printf("Esperando aviso")
			ma := <-msgsAviso
			if len(ma.Body) > 0 {
				log.Printf("Aviso: %s", ma.Body)
				ma.Ack(false)
			}
			// Send a message to the gangster queue
			err = ch.Publish(
				"",         // exchange
				"gangster", // routing key
				false,      // mandatory
				false,      // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte("Entra gangster"),
				})
			failOnError(err, "Failed to publish a message")

			d := <-msgs
			log.Printf("Received a message: %s", d.Body)
			d.Ack(false)
			numshuhisAComer--
			rand.Seed(time.Now().UnixNano())
			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			if numshuhisAComer == 0 {
				done <- Empty{}
				break
			}

			//TODO: Borrar la cola de mensajes
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-done
}