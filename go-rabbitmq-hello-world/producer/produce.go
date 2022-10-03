package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	must(err)
	defer conn.Close()

	ch, err1 := conn.Channel()
	must(err1)
	defer ch.Close()

	q, err2 := ch.QueueDeclare(
		"messagesQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	must(err2)

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
	must(err3)

	log.Println("Message published to the queue successfully :)")

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
