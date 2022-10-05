package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func HandleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	conn, err0 := amqp.Dial("amqp://guest:guest@localhost:5672/")
	HandleError(err0, "Cannot connect to AMQP")
	defer conn.Close()

	ch, err1 := conn.Channel()
	HandleError(err1, "Cannot create a amqpChannel")
	defer ch.Close()

	err2 := ch.ExchangeDeclare(
		"logs",
		"fanout",
		true, // durable
		false,
		false,
		false,
		nil,
	)
	HandleError(err2, "Failed to declare an exchange")

	q, err3 := ch.QueueDeclare(
		"",
		false,
		false,
		true, //exclusive
		false,
		nil,
	)
	HandleError(err3, "Could not declare a temp queue")

	err4 := ch.QueueBind(
		q.Name,
		"",
		"logs",
		false,
		nil,
	)
	HandleError(err4, "Failed to bind a queue")

	msgs, err5 := ch.Consume(
		q.Name,
		"",
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)
	HandleError(err5, "Could not register consumer")

	forever := make(chan string)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
