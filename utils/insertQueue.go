package utils

import (
	"encoding/json"
	"log"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"github.com/streadway/amqp"
)

func StartInsertQueue(conn *amqp.Connection) {

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	err = ch.ExchangeDeclare("notifications", "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare the exchange: %v", err)
	}

	queue, err := ch.QueueDeclare(
		"insert-notification",
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
		"insert",
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
			log.Printf("Received a message: %v", msg.AppId)
			var notification models.TransactionNotificationVariables
			err := json.Unmarshal(msg.Body, &notification)
			if err != nil {
				log.Println(err)

			} else {
				InsertUserAndNotification(notification)
			}
		}
		defer ch.Close()
	}()
}
