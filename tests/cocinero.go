package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
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
	numRanT1 := rand.Intn(TO_PRODUCE)
	if numRanT1 != 0 {
		log.Printf("Preparando %d de %s", numRanT1, SHUSHI_TYPES[0])
	}
	restantes := TO_PRODUCE - numRanT1
	numRanT2 := rand.Intn(restantes)
	if numRanT2 != 0 {
		log.Printf("Preparando %d de %s", numRanT2, SHUSHI_TYPES[1])
	}
	restantes = restantes - numRanT2
	if restantes != 0 {
		log.Printf("Preparando %d de %s", restantes, SHUSHI_TYPES[2])
	}

	var shushi = ""
	for i := 0; i < TO_PRODUCE; i++ {
		if numRanT1 != 0 {
			shushi = SHUSHI_TYPES[0]
			numRanT1--
		} else if numRanT2 != 0 {
			shushi = SHUSHI_TYPES[1]
			numRanT2--
		} else {
			shushi = SHUSHI_TYPES[2]
		}
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
		// Esperar un tiempo aleatorio entre 1 y 2 segundos
		//rand.Seed(time.Now().UnixNano())
		//time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
	}
}

func avisar(ch *amqp.Channel, ctx context.Context, aviso amqp.Queue, err error) {
	//	for i := 0; i < TO_PRODUCE; i++ {
	//body := fmt.Sprintf("Aviso: %d", i)
	body := fmt.Sprintf("%d", TO_PRODUCE)
	err = ch.PublishWithContext(ctx,
		"",         // exchange
		aviso.Name, // routing key
		false,      // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	//}
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

	aviso, err := ch.QueueDeclare(
		"avisos", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	producir()

	enviar(ch, ctx, q, err)

	avisar(ch, ctx, aviso, err)
	// En vez de enviar N permisos enviar un int con N valor y que cada cliente modifique el valor
}
