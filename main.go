package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Kofi-D-Boateng/legacynotifications/controllers"
	"github.com/Kofi-D-Boateng/legacynotifications/database"
	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"github.com/Kofi-D-Boateng/legacynotifications/router"
	"github.com/Kofi-D-Boateng/legacynotifications/utils"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	os.Setenv("GO_ENV", "dev")
	env := os.Getenv("GO_ENV")
	if env == "dev" {
		_, file, _, ok := runtime.Caller(0)
		basePath := filepath.Dir(file)
		fmt.Println(file)
		fmt.Println(basePath)

		if !ok {
			log.Fatalf("Unable to find file path: %v", file)
		}

		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	}

	// GRABBING DB INFO
	var dns string = os.Getenv("MONGO_URI")
	var dbName string = os.Getenv("DB_NAME")
	database.Database.UserCollection = os.Getenv("USERS_COLLECTION")
	database.Database.CustomerServiceCollection = os.Getenv("CUSTOMER_SERVICE_COLLECTION")

	// Set options
	clientOptions := options.Client().ApplyURI(dns)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("MongoDB connected to DB: %v \n", dbName)
	// Get DB
	database.Database.Db = client.Database(dbName)
}

func main() {

	// Setting up rabbitmq
	queueConnStr := os.Getenv("RABBITMQ_CONN")
	if queueConnStr == "" {
		log.Fatalln("[ERROR]: RabbitMQ Connection String is missing")
	}
	conn, err := amqp.Dial(queueConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
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
			database.MarkMessageAsRead(markMessage)
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
				database.InsertUserAndNotification(notification)
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
				utils.SendConfirmationEmail(emailAttributes)
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
				utils.SendMailingListConfirmation(email)
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
				controllers.CustomerServiceRequest(request)
			}
			log.Printf("Received a message: %v", request)
		}
	}()

	r := router.Router()
	port := os.Getenv("PORT")

	srv := &http.Server{
		Handler:      r,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Server listening at port%v \n", port)

	log.Fatal(srv.ListenAndServe())
}
