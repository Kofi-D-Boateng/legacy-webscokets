package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Database
var UserCollection string;
var CustomerService string;


func init(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// GRABBING DB INFO
	var dns string = os.Getenv("MONGO_URI")
	var dbName string = os.Getenv("DB_NAME")
	UserCollection = os.Getenv("USERS_COLLECTION")
	CustomerService = os.Getenv("CUSTOMER_SERVICE_COLLECTION")

	// Set options
	clientOptions := options.Client().ApplyURI(dns)
	client, err := mongo.Connect(context.TODO(),clientOptions)

	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("MongoDB connected")
	// Get DB
	Db = client.Database(dbName)
}

func FindAUser(email string) models.User{
	var result models.User
	users := Db.Collection(UserCollection)
	filter := bson.M{"email":email}
	err := users.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func MarkMessageAsRead(request models.MarkMessage) bool {
	fmt.Printf("\n email: %v", request.Email)
	fmt.Printf("\n email: %v", request.MsgID)
	return true
}

func InsertUserAndNotification(variables struct {Email					string `json:"email"`
		Receiver				string `json:"receiver" `
		ReceiverEmail			string `json:"receiverEmail"`
		Sender					string `json:"sender"`
		IsReceiverInDatabase 	bool `json:"isReceiverInDatabase"`
		DateOfTransaction 		string `json:"dateOfTransaction"`
		Type 					string `json:"type"`
		Amount 					float64 `json:"amount"`}) bool {

	var transaction models.Transaction
	var sender models.User
	var receiver models.User
	receiverEmailFilter := bson.M{"email": variables.ReceiverEmail}
	senderEmailFilter := bson.M{"email": variables.Email}
	users := Db.Collection(UserCollection)

	transaction.Amount = variables.Amount
	transaction.DateOfTransaction = variables.DateOfTransaction
	transaction.Receiver = variables.Receiver
	transaction.ReceiverEmail = variables.ReceiverEmail
	transaction.Sender = variables.Sender
	transaction.Type = variables.Type

	fmt.Println(transaction)

	// FIND PERSONNEL
	errOne := users.FindOne(context.Background(), receiverEmailFilter).Decode(&receiver)
	errTwo := users.FindOne(context.Background(),senderEmailFilter).Decode(&sender)
	if errOne != nil {
		log.Fatal(errOne)
	}
	if errTwo != nil {
		log.Fatal(errTwo)
	}

	// BUSINESS LOGIC

	// CREATE NOTIFICATION FOR IN-HOUSE CUSTOMERS
	if variables.IsReceiverInDatabase && receiver.Email != variables.ReceiverEmail {
		receiver.Email = variables.ReceiverEmail
		receiver.Notifications = []models.Transaction{}
		_, err := users.UpdateOne(context.Background(),receiverEmailFilter,receiver)
		if err != nil {
			log.Fatal(err)
		}
	}

	if variables.Email != sender.Email {
		sender.Email = variables.Email
		sender.Notifications = []models.Transaction{}
		_, err := users.UpdateOne(context.Background(), senderEmailFilter, sender)
		if err != nil {
			log.Fatal(err)
		}
		return true
	}

	// UPDATING IN-HOUSE RECEIVER
	if variables.IsReceiverInDatabase && receiver.Email == variables.ReceiverEmail {
		receiver.Notifications = append(receiver.Notifications, transaction)
		_,err := users.UpdateOne(context.Background(), receiverEmailFilter, receiver)

		if err != nil {
			log.Fatal(err)
		}
	}

	sender.Notifications = append(sender.Notifications, transaction)
	_,err := users.UpdateOne(context.Background(), senderEmailFilter, sender)

	if err != nil {
		log.Fatal(err)
	}

	return true
}