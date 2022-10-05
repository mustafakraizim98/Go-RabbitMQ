package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type User struct {
	Username string    `json:"username,omitempty" binding:"required"`
	Password string    `json:"password,omitempty" binding:"required"`
	LoggedAt time.Time `json:"logged_at,omitempty"`
}

func HandleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func JsonMarshal(user User) []byte {
	json, err := json.Marshal(user)
	HandleError(err, "Could not marshaling the json object")
	return json
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := User{
		Username: "mustafa98",
		Password: "0123456789",
		LoggedAt: time.Now(),
	}

	body := JsonMarshal(user)

	err3 := ch.PublishWithContext(
		ctx,
		"logs",
		"",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		},
	)
	HandleError(err3, "Error publishing message")

	log.Println("Message published to the Exchange successfully :)")
}
