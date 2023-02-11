package utils

import (
	"encoding/json"
	"log"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"github.com/streadway/amqp"
)

func StartUpdateQueue(conn *amqp.Connection) {

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	err = ch.ExchangeDeclare("notifications", "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare the exchange: %v", err)
	}

	queue, decErr := ch.QueueDeclare(
		"update-notifications",
		true,
		false,
		false,
		false,
		nil,
	)

	if decErr != nil {
		log.Fatalf("Failed to bind queue to exchange: %v", err)
	}

	err = ch.QueueBind(
		queue.Name,
		"update",
		"notifications",
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Failed to bind queue to exchange: %v", err)
	}

	go func() {
		msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
		if err != nil {
			log.Fatalf("Error consuming from %v. %v", queue.Name, err)
		}
		for msg := range msgs {
			msg.Ack(false)
			var markMessage models.MarkMessage
			log.Printf("Received a message: %v", msg)
			err = json.Unmarshal([]byte(msg.Body), &markMessage)
			if err != nil {
				log.Printf("Error unmarshalling message: %s", err)
				continue
			}
			log.Printf("Received a message: %v", markMessage)
			MarkMessageAsRead(markMessage)
		}
		defer ch.Close()
	}()
}
