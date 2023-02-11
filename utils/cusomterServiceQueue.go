package utils

import (
	"encoding/json"
	"log"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"github.com/streadway/amqp"
)

func StartCustomerServiceQueue(conn *amqp.Connection) {

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	err = ch.ExchangeDeclare("notifications", "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare the exchange: %v", err)
	}

	queue, err := ch.QueueDeclare(
		"customer-service",
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
		"cust-serv",
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
			var request models.CustomerServiceMessage
			err = json.Unmarshal([]byte(msg.Body), &request)
			if err != nil {
				log.Printf("Error unmarshalling message: %s", err)

			} else {
				if !EmailExpr.Match([]byte(request.Emailer)) {
					log.Fatalf("Error with email: %s", request.Emailer)
				}

				if Accounts.Match([]byte(request.Topic)) {
					SendToAccountDept(request)
				}

				if Billing.Match([]byte(request.Topic)) {
					SendToBillingDept(request)
				}

				SendToOther(request)
			}
			log.Printf("Received a message: %v", request)
		}
		defer ch.Close()
	}()
}
