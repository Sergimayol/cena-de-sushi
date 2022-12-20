package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
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
		0,     // prefetch count
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

	rand.Seed(time.Now().UnixNano())
	numshuhisAComer := (rand.Intn(11) + 1)

	fmt.Printf("Bon vespre vinc a sopar de shushi\nAvui menjaré %d peçes de shushi\n",
		numshuhisAComer)

	done := make(chan Empty, 1)

	// Receive messages from the queue
	go func() {
		// Get the messages from the queue
		for {

			ma := <-msgsAviso
			if len(ma.Body) > 0 {
				ma.Ack(false)
			}

			// Convert string to int
			i, err := strconv.Atoi(string(ma.Body))
			failOnError(err, "Failed to convert string to int")
			if i > 0 {
				i--
				err = ch.Publish(
					"",         // exchange
					aviso.Name, // routing key
					false,      // mandatory
					false,      // immediate
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(fmt.Sprintf("%d", i)),
					})
				failOnError(err, "Failed to publish a message")
				d := <-msgs
				log.Printf("Ha agafat un %s\nAl plat hi ha %d peces", d.Body, i)
				d.Ack(false)
				numshuhisAComer--
				// Esperar un segundo
				time.Sleep(1 * time.Second)
				if numshuhisAComer == 0 {
					done <- Empty{}
					break
				}
			}

		}
		fmt.Println("Ja estic ple. Bona nit")
	}()

	// Esperar a que acabe
	<-done
}
