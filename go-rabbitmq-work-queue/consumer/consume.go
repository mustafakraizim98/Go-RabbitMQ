package main

import (
	"encoding/json"
	"fmt"
	"log"
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

	err3 := ch.Qos(
		1,
		0,
		false,
	)
	handleError(err3, "Could not configure QoS")

	msgs, err4 := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err4, "Could not register consumer")

	forever := make(chan string)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			task := Task{}

			err := json.Unmarshal(d.Body, &task)
			handleError(err, "Error decoding JSON")

			numberA := task.NumberA
			numberB := task.NumberB
			result := numberA + numberB
			fmt.Printf("NumberA: %d, \nNumberB: %d, \nResult: %d\n", numberA, numberB, result)

			t := time.Duration(q.Messages)
			time.Sleep(t * time.Second)
			fmt.Println("Done")

			if err := d.Ack(false); err != nil {
				handleError(err, "Error acknowledging message")
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
