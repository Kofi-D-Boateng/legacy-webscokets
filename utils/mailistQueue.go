package utils

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func StartMaillistQueue(conn *amqp.Connection) {

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	err = ch.ExchangeDeclare("notifications", "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare the exchange: %v", err)
	}

	queue, err := ch.QueueDeclare(
		"join-maillist",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare the queue: %v", err)
	}

	err = ch.QueueBind(
		queue.Name,
		"maillist",
		"notifications",
		false,
		nil,
	)
	fmt.Printf("QUEUE NAME & CONSUMERS: %v & %v\n", queue.Name, queue.Consumers)
	if err != nil {
		log.Fatalf("Failed to bind queue to exchange: %v", err)
	}

	go func() {
		msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
		if err != nil {
			log.Fatalf("Error consuming from %v. %v", queue.Name, err)
		}
		for msg := range msgs {
			log.Printf("Received a message in %v\n", queue.Name)
			var maillist struct {
				Email string `json:"email"`
			}
			err = json.Unmarshal([]byte(msg.Body), &maillist.Email)
			if err != nil {
				log.Printf("Error unmarshalling message: %s", err)

			} else {
				fmt.Printf("EMAIL: %v\n", maillist.Email)
				SendMailingListConfirmation(maillist.Email)
			}
		}
		defer ch.Close()
	}()

}
