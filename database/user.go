package database

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func FindAUser(email string) models.User {
	var result models.User
	userCollection := Database.Db.Collection(Database.UserCollection)
	filter := bson.M{"email": email}
	err := userCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Print(err)
	}
	return result
}

func MarkMessageAsRead(request models.MarkMessage) models.User {
	fmt.Printf("\n email: %v \n", request.Email)
	fmt.Printf("\n id: %v \n", request.MsgID)

	id,err := primitive.ObjectIDFromHex(request.MsgID)
	if err != nil {
		fmt.Printf("Invalid hex string: %v \n",err)
	}
	userCollectionPointer := Database.Db.Collection(Database.UserCollection)
	filter := bson.M{"email": request.Email}
	update := bson.M{"$set": bson.M{"notifications.$[element].read":true }}
	arrayFilterOptions := options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
        Filters: []interface{}{bson.M{"element._id": id}},
    })
	
	result := userCollectionPointer.FindOneAndUpdate(context.Background(), filter, update, arrayFilterOptions)

	if result == nil {
		fmt.Printf("Error updating document for:%v",request.Email)
	}

	customer := FindAUser(request.Email)

	return customer
}

func InsertUserAndNotification(variables models.TransactionNotificationVariables) int {


	var status struct {
		isReceiverUpdated 	bool
		isSenderUpdated 	bool
	}

	fmt.Printf("RECEIVER STATUS: %v \n SENDER STATUS: %v \n", status.isReceiverUpdated,status.isSenderUpdated)

	var transaction models.Transaction
	var sender models.User
	var receiver models.User
	receiverEmailFilter := bson.M{"email": variables.ReceiverEmail}
	senderEmailFilter := bson.M{"email": variables.Email}
	userCollectionPointer := Database.Db.Collection(Database.UserCollection)

	transaction.ID = primitive.NewObjectID()
	transaction.Amount = variables.Amount
	transaction.Date = variables.DateOfTransaction
	transaction.Receiver = variables.Receiver
	transaction.Sender = variables.Sender
	transaction.Read = false

	fmt.Println(transaction)

	// BUSINESS LOGIC
  
	// FIND PERSONNEL
	errOne := userCollectionPointer.FindOne(context.Background(), receiverEmailFilter).Decode(&receiver)
	errTwo := userCollectionPointer.FindOne(context.Background(), senderEmailFilter).Decode(&sender)

	// USER HAS NOT CREATE INSTANCES FOR PURCHASES YET.
	if errOne != nil {
		log.Printf("COULD NOT FIND RECIPIENT: %v in database, Attempting to create notifications if user is in main db.... \n",variables.Receiver)
		receiver.ID = primitive.NewObjectID()
		if variables.IsReceiverInDatabase {
			receiver.Email = variables.ReceiverEmail
			receiver.Notifications = []models.Transaction{transaction}

			_, errForReceiver := userCollectionPointer.InsertOne(context.Background(), receiver)

			if errForReceiver  != nil {
				log.Printf("ERROR INSERTING DOCUMENTS FOR %v \n", receiver.Email)
				return http.StatusInternalServerError
			}
		}
		status.isReceiverUpdated = true
	}


	// USER HAS NOT CREATE INSTANCES FOR PURCHASES YET.
	if errTwo != nil {
		log.Printf("COULD NOT FIND TRASFERER: %v in database, Attempting to create notifications \n",variables.Sender)
		sender.ID = primitive.NewObjectID()
		sender.Email = variables.Email

		sender.Notifications = []models.Transaction{transaction}

		_, errForSender := userCollectionPointer.InsertOne(context.Background(), sender)

		if errForSender  != nil {
			log.Printf("ERROR INSERTING DOCUMENTS FOR, %v \n", sender.Email)
			return http.StatusInternalServerError
		}
		status.isSenderUpdated = true
	}

	if status.isReceiverUpdated && status.isSenderUpdated {
		return http.StatusOK
	}


	// UPDATING IN-HOUSE RECEIVER & TRANSFERER

	if receiver.Email == variables.ReceiverEmail && !status.isReceiverUpdated {
		receiver.Notifications = append(receiver.Notifications, transaction)

		receiverUpdate := bson.M{"$set": bson.M{"notifications" : receiver.Notifications}}

		resultTwo := userCollectionPointer.FindOneAndUpdate(context.Background(), receiverEmailFilter, receiverUpdate)

		if resultTwo.Err() == mongo.ErrNoDocuments {
			log.Printf("ERROR WITH FINDING AND UPDATING FOR %v: %v \n", receiver.Email ,resultTwo)
			return http.StatusInternalServerError
		}
		status.isReceiverUpdated = true
	}

	if status.isReceiverUpdated && status.isSenderUpdated {
		status.isReceiverUpdated = false
		status.isSenderUpdated = false
		return http.StatusOK
	}


	if variables.Email == sender.Email && !status.isSenderUpdated {

		sender.Notifications = append(sender.Notifications, transaction)
		senderUpdate := bson.M{"$set": bson.M{"notifications" : sender.Notifications}}

		resultOne := userCollectionPointer.FindOneAndUpdate(context.Background(), senderEmailFilter, senderUpdate)

		if resultOne.Err() == mongo.ErrNoDocuments {
			log.Printf("ERROR WITH FINDING AND UPDATING FOR %v: %v \n", sender.Email ,resultOne)
			return http.StatusInternalServerError
		}
	}

	status.isReceiverUpdated = false
	status.isSenderUpdated = false
	return http.StatusOK
}