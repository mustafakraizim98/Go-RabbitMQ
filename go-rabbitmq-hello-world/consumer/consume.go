package main

import (
	"log"

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

	msgs, err3 := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	must(err3)

	forever := make(chan string)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
