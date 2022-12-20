package main

import (
	"fmt"
	"log"
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

	if len(msgs) > 0 {
	}

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

	fmt.Printf("Bon vespre vinc a sopar de shushi\nHo vull tot!\n")

	done := make(chan Empty, 1)

	// Receive messages from the queue
	go func() {
		// Get the messages from the queue
		ma := <-msgsAviso
		ma.Ack(false)

		i, err := strconv.Atoi(string(ma.Body))
		failOnError(err, "Failed to convert string to int")
		fmt.Printf("Agafa totes les peces. En total %d\n", i)

		/*err = ch.Publish(
			"",         // exchange
			aviso.Name, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(fmt.Sprintf("%d", 0)),
			})
		failOnError(err, "Failed to publish a message")*/

		// Clean all the messages from the queue
		_, err = ch.QueuePurge(q.Name, false)
		_, err = ch.QueueDelete(q.Name, false, false, false)
		failOnError(err, "Failed to purge the queue")

		fmt.Println("Romp el plat")

		// Esperar 1 segundo
		time.Sleep(1 * time.Second)

		fmt.Println("Men vaig")

		done <- Empty{}

	}()

	<-done
}
