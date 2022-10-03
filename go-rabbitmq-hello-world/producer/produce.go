package main

import (
	"context"
	"log"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello, World!"
	err3 := ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	handleError(err3, "Error publishing message")

	log.Println("Message published to the queue successfully :)")

}
