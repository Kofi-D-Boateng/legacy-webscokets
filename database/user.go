package database

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"go.mongodb.org/mongo-driver/bson"
)


func FindAUser(email string) models.User {
	var result models.User
	users := Db.Collection(UserCollection)
	filter := bson.M{"email": email}
	err := users.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Print(err)
	}
	return result
}

func MarkMessageAsRead(request models.MarkMessage) int {
	fmt.Printf("\n email: %v", request.Email)
	fmt.Printf("\n id: %v", request.MsgID)

	var customer models.User

	user := Db.Collection(UserCollection)
	filter := bson.M{"email": request.Email}
	err := user.FindOne(context.Background(), filter).Decode(&customer)

	if err != nil {
		log.Print(err)
		return http.StatusInternalServerError
	}

	for _, noti := range customer.Notifications {
		if noti.ID.Hex() == request.MsgID {
			noti.Read = true

		} else {
			continue
		}
	}

	_, updateErr := user.UpdateOne(context.Background(), filter, customer)

	if updateErr != nil {
		log.Print(updateErr)
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func InsertUserAndNotification(variables struct {
	Email                string  `json:"email"`
	Receiver             string  `json:"receiver" `
	ReceiverEmail        string  `json:"receiverEmail"`
	Sender               string  `json:"sender"`
	IsReceiverInDatabase bool    `json:"isReceiverInDatabase"`
	DateOfTransaction    string  `json:"dateOfTransaction"`
	Type                 string  `json:"type"`
	Amount               float64 `json:"amount"`
}) int {

	var transaction models.Transaction
	var sender models.User
	var receiver models.User
	receiverEmailFilter := bson.M{"email": variables.ReceiverEmail}
	senderEmailFilter := bson.M{"email": variables.Email}
	users := Db.Collection(UserCollection)

	transaction.Amount = variables.Amount
	transaction.Date = variables.DateOfTransaction
	transaction.Receiver = variables.Receiver
	transaction.Sender = variables.Sender
	transaction.Read = false

	fmt.Println(transaction)

	// FIND PERSONNEL
	errOne := users.FindOne(context.Background(), receiverEmailFilter).Decode(&receiver)
	errTwo := users.FindOne(context.Background(), senderEmailFilter).Decode(&sender)
	if errOne != nil {
		log.Print(errOne)
		return http.StatusInternalServerError
	}
	if errTwo != nil {
		log.Print(errTwo)
		return http.StatusInternalServerError
	}

	// BUSINESS LOGIC

	// CREATE NOTIFICATION FOR IN-HOUSE CUSTOMERS
	if variables.IsReceiverInDatabase && receiver.Email != variables.ReceiverEmail {
		receiver.Email = variables.ReceiverEmail
		receiver.Notifications = []models.Transaction{}
		_, err := users.UpdateOne(context.Background(), receiverEmailFilter, receiver)
		if err != nil {
			log.Print(err)
			return http.StatusInternalServerError
		}
	}

	if variables.Email != sender.Email {
		sender.Email = variables.Email
		sender.Notifications = []models.Transaction{}
		_, err := users.UpdateOne(context.Background(), senderEmailFilter, sender)
		if err != nil {
			log.Print(err)
			return http.StatusInternalServerError
		}
		return http.StatusOK
	}

	// UPDATING IN-HOUSE RECEIVER
	if variables.IsReceiverInDatabase && receiver.Email == variables.ReceiverEmail {
		receiver.Notifications = append(receiver.Notifications, transaction)
		_, err := users.UpdateOne(context.Background(), receiverEmailFilter, receiver)

		if err != nil {
			log.Print(err)
			return http.StatusInternalServerError
		}
	}

	sender.Notifications = append(sender.Notifications, transaction)
	_, err := users.UpdateOne(context.Background(), senderEmailFilter, sender)

	if err != nil {
		log.Print(err)
		return http.StatusInternalServerError
	}

	return http.StatusOK
}