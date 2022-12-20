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

func comer(ch *amqp.Channel, msgs <-chan amqp.Delivery, msgsAviso <-chan amqp.Delivery, done chan<- Empty, aviso amqp.Queue, numshuhisAComer int) {
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
}

func declararCola(ch *amqp.Channel, nombre string) amqp.Queue {
	q, err := ch.QueueDeclare(
		nombre, // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return q
}

func crearCanalConsumir(ch *amqp.Channel, q amqp.Queue) <-chan amqp.Delivery {
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
	return msgs
}

func main() {
	conn, err := amqp.Dial(RABBITMQ_HOST)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q := declararCola(ch, "task_queue")
	aviso := declararCola(ch, "avisos")

	err = ch.Qos(
		0,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs := crearCanalConsumir(ch, q)
	msgsAviso := crearCanalConsumir(ch, aviso)

	rand.Seed(time.Now().UnixNano())
	numshuhisAComer := (rand.Intn(11) + 1)

	fmt.Printf("Bon vespre vinc a sopar de shushi\nAvui menjaré %d peçes de shushi\n",
		numshuhisAComer)

	done := make(chan Empty, 1)

	go comer(ch, msgs, msgsAviso, done, aviso, numshuhisAComer)

	// Esperar a que acabe
	<-done
}
