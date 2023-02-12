package utils

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"github.com/streadway/amqp"
)

func StartVerificationQueue(conn *amqp.Connection) {

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	err = ch.ExchangeDeclare("notifications", "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare the exchange: %v", err)
	}

	queue, err := ch.QueueDeclare(
		"verification-email",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare the queue: %v", err)
	}
	fmt.Printf("QUEUE NAME & CONSUMERS: %v & %v\n", queue.Name, queue.Consumers)
	err = ch.QueueBind(
		queue.Name,
		"verification",
		"notifications",
		false,
		nil,
	)
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
			var emailAttributes models.EmailAttributes
			err = json.Unmarshal([]byte(msg.Body), &emailAttributes)
			if err != nil {
				log.Printf("Error unmarshalling message: %s", err)

			} else {
				SendConfirmationEmail(emailAttributes)
			}
			log.Printf("Received a message: %v", emailAttributes)
		}
		defer ch.Close()
	}()
}
