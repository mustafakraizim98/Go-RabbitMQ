package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	conn, err0 := amqp.Dial("amqp://guest:guest@localhost:5672/")
	handleError(err0, "Cannot connect to AMQP")
	defer conn.Close()

	ch, err1 := conn.Channel()
	handleError(err1, "Cannot create a amqpChannel")
	defer ch.Close()

	q, err2 := ch.QueueDeclare(
		"messagesQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err2, "Could not declare `messagesQueue` queue")

	msgs, err3 := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	handleError(err3, "Could not register consumer")

	forever := make(chan string)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
