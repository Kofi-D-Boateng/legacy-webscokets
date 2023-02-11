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
		"maillist-verification",
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
			var email string
			err = json.Unmarshal([]byte(msg.Body), &email)
			if err != nil {
				log.Printf("Error unmarshalling message: %s", err)

			} else {
				fmt.Println(email)
				SendMailingListConfirmation(email)
			}
			log.Printf("Received a message: %v", email)
		}
		defer ch.Close()
	}()

}
