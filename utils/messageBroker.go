package utils

import (
	"encoding/json"
	"log"

	"github.com/Kofi-D-Boateng/legacynotifications/models"

	"github.com/streadway/amqp"
)

var RabbitMQ *amqp.Connection

func ConnectMessageBroker(RabbitMQ *amqp.Connection) {

	ch, err := RabbitMQ.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare("notifications", "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare the exchange: %v", err)
	}

	queue1, err := ch.QueueDeclare(
		"update-notifications",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare the queue: %v", err)
	}

	queue2, err := ch.QueueDeclare(
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

	queue3, err := ch.QueueDeclare(
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

	queue4, err := ch.QueueDeclare(
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

	queue5, err := ch.QueueDeclare(
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
		queue1.Name,
		"update",
		"notifications",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue to exchange: %v", err)
	}

	err = ch.QueueBind(
		queue2.Name,
		"insert",
		"notifications",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue to exchange: %v", err)
	}

	err = ch.QueueBind(
		queue3.Name,
		"verification",
		"notifications",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue to exchange: %v", err)
	}

	err = ch.QueueBind(
		queue4.Name,
		"mailist",
		"notifications",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue to exchange: %v", err)
	}

	err = ch.QueueBind(
		queue5.Name,
		"cust-serv",
		"notifications",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue to exchange: %v", err)
	}

	go func() {
		q1Msgs, err := ch.Consume(queue1.Name, "", true, false, false, false, nil)
		if err != nil {
			log.Fatalf("Error consuming from %v. %v", queue1.Name, err)
		}
		for msg := range q1Msgs {
			var markMessage models.MarkMessage
			err = json.Unmarshal(msg.Body, &markMessage)
			if err != nil {
				log.Printf("Error unmarshalling message: %s", err)
				continue
			}
			log.Printf("Received a message: %v", markMessage)
			MarkMessageAsRead(markMessage)
		}
	}()

	go func() {
		q2Msgs, err := ch.Consume(queue2.Name, "", true, false, false, false, nil)
		if err != nil {
			log.Fatalf("Error consuming from %v. %v", queue2.Name, err)
		}
		for msg := range q2Msgs {
			var notification models.TransactionNotificationVariables
			err = json.Unmarshal(msg.Body, &notification)
			if err != nil {
				log.Printf("Error unmarshalling message: %s", err)

			} else {
				InsertUserAndNotification(notification)
			}
			log.Printf("Received a message: %v", notification)
		}
	}()

	go func() {
		q3Msgs, err := ch.Consume(queue3.Name, "", true, false, false, false, nil)
		if err != nil {
			log.Fatalf("Error consuming from %v. %v", queue3.Name, err)
		}
		for msg := range q3Msgs {
			var emailAttributes models.EmailAttributes
			err = json.Unmarshal(msg.Body, &emailAttributes)
			if err != nil {
				log.Printf("Error unmarshalling message: %s", err)

			} else {
				SendConfirmationEmail(emailAttributes)
			}
			log.Printf("Received a message: %v", emailAttributes)
		}
	}()

	go func() {
		q4Msgs, err := ch.Consume(queue4.Name, "", true, false, false, false, nil)
		if err != nil {
			log.Fatalf("Error consuming from %v. %v", queue4.Name, err)
		}
		for msg := range q4Msgs {
			var email string
			err = json.Unmarshal(msg.Body, &email)
			if err != nil {
				log.Printf("Error unmarshalling message: %s", err)

			} else {
				SendMailingListConfirmation(email)
			}
			log.Printf("Received a message: %v", email)
		}
	}()

	go func() {
		q5Msgs, err := ch.Consume(queue5.Name, "", true, false, false, false, nil)
		if err != nil {
			log.Fatalf("Error consuming from %v. %v", queue5.Name, err)
		}
		for msg := range q5Msgs {
			var request models.CustomerServiceMessage
			err = json.Unmarshal(msg.Body, &request)
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
	}()
}
