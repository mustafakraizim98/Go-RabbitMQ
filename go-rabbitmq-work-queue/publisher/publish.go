package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	amqp "github.com/streadway/amqp"
)

type Task struct {
	NumberA int
	NumberB int
}

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
		"workQueue",
		true,
		false,
		false,
		false,
		nil,
	)
	handleError(err2, "Could not declare `workQueue` queue")

	// seed to get different result every time
	rand.Seed(time.Now().UnixNano())

	task := Task{NumberA: rand.Intn(10), NumberB: rand.Intn(30)}
	body, err3 := json.Marshal(task)
	handleError(err3, "Error encoding JSON")

	err4 := ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		},
	)
	handleError(err4, "Error publishing message")

	log.Println("Message published to the queue successfully :)")

}
