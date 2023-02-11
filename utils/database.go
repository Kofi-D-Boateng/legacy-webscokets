package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	EmailExpr *regexp.Regexp
	Accounts  *regexp.Regexp
	Billing   *regexp.Regexp
	Database  *mongo.Database
)

func ConnectDatabase(uri string, dbName string) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}

	pingErr := client.Ping(context.Background(), nil)

	if pingErr != nil {
		panic(pingErr)
	}

	Database = client.Database(dbName)
}

func FindAUser(email string) models.User {
	var result models.User
	userCollection := Database.Collection(os.Getenv("USER_COLLECTION"))
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

	id, err := primitive.ObjectIDFromHex(request.MsgID)
	if err != nil {
		fmt.Printf("Invalid hex string: %v \n", err)
	}
	userCollectionPointer := Database.Collection(os.Getenv("USER_COLLECTION"))
	filter := bson.M{"email": request.Email}
	update := bson.M{"$set": bson.M{"notifications.$[element].read": true}}
	arrayFilterOptions := options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"element._id": id}},
	})

	result := userCollectionPointer.FindOneAndUpdate(context.Background(), filter, update, arrayFilterOptions)

	if result == nil {
		fmt.Printf("Error updating document for:%v", request.Email)
	}

	customer := FindAUser(request.Email)

	return customer
}

func InsertUserAndNotification(request models.TransactionNotificationVariables) {

	var status struct {
		isReceiverUpdated bool
		isSenderUpdated   bool
	}

	var transaction models.Transaction
	var sender models.User
	var receiver models.User

	receiverEmailFilter := bson.M{"email": request.ReceiverEmail}
	senderEmailFilter := bson.M{"email": request.Email}
	userCollectionPointer := Database.Collection(os.Getenv("USER_COLLECTION"))

	dot := time.Date(request.DateOfTransaction[0], time.Month(request.DateOfTransaction[1]), request.DateOfTransaction[2], request.DateOfTransaction[3], request.DateOfTransaction[4], request.DateOfTransaction[5], request.DateOfTransaction[6], time.UTC)
	transaction.ID = primitive.NewObjectID()
	transaction.Amount = request.Amount
	transaction.Date = dot.Format("2006-01-02T15:04:05.999")
	transaction.Receiver = request.Receiver
	transaction.Sender = request.Sender
	transaction.Read = false

	fmt.Println(transaction)

	// BUSINESS LOGIC

	// FIND PERSONNEL
	errOne := userCollectionPointer.FindOne(context.Background(), receiverEmailFilter).Decode(&receiver)
	errTwo := userCollectionPointer.FindOne(context.Background(), senderEmailFilter).Decode(&sender)

	// USER HAS NOT CREATE INSTANCES FOR PURCHASES YET.
	if errOne != nil {
		log.Printf("COULD NOT FIND RECIPIENT: %v in database, Attempting to create notifications if user is in main db.... \n", request.Receiver)
		receiver.ID = primitive.NewObjectID()
		if request.IsReceiverInDatabase {
			receiver.Email = request.ReceiverEmail
			receiver.Notifications = []models.Transaction{transaction}

			_, errForReceiver := userCollectionPointer.InsertOne(context.Background(), receiver)

			if errForReceiver != nil {
				log.Printf("ERROR INSERTING DOCUMENTS FOR %v \n", receiver.Email)
				return
			}
		}
		status.isReceiverUpdated = true
	}

	// USER HAS NOT CREATE INSTANCES FOR PURCHASES YET.
	if errTwo != nil {
		log.Printf("COULD NOT FIND TRASFERER: %v in database, Attempting to create notifications \n", request.Sender)
		sender.ID = primitive.NewObjectID()
		sender.Email = request.Email

		sender.Notifications = []models.Transaction{transaction}

		_, errForSender := userCollectionPointer.InsertOne(context.Background(), sender)

		if errForSender != nil {
			log.Printf("ERROR INSERTING DOCUMENTS FOR, %v \n", sender.Email)
			return
		}
		status.isSenderUpdated = true
	}

	if status.isReceiverUpdated && status.isSenderUpdated {
		return
	}

	// UPDATING IN-HOUSE RECEIVER & TRANSFERER

	if receiver.Email == request.ReceiverEmail && !status.isReceiverUpdated {
		receiver.Notifications = append(receiver.Notifications, transaction)

		receiverUpdate := bson.M{"$set": bson.M{"notifications": receiver.Notifications}}

		resultTwo := userCollectionPointer.FindOneAndUpdate(context.Background(), receiverEmailFilter, receiverUpdate)

		if resultTwo.Err() == mongo.ErrNoDocuments {
			log.Printf("ERROR WITH FINDING AND UPDATING FOR %v: %v \n", receiver.Email, resultTwo)
			return
		}
		status.isReceiverUpdated = true
	}

	if status.isReceiverUpdated && status.isSenderUpdated {
		status.isReceiverUpdated = false
		status.isSenderUpdated = false
		return
	}

	if request.Email == sender.Email && !status.isSenderUpdated {

		sender.Notifications = append(sender.Notifications, transaction)
		senderUpdate := bson.M{"$set": bson.M{"notifications": sender.Notifications}}

		resultOne := userCollectionPointer.FindOneAndUpdate(context.Background(), senderEmailFilter, senderUpdate)

		if resultOne.Err() == mongo.ErrNoDocuments {
			log.Printf("ERROR WITH FINDING AND UPDATING FOR %v: %v \n", sender.Email, resultOne)
			return
		}
	}

	status.isReceiverUpdated = false
	status.isSenderUpdated = false
}

func SendToOther(details models.CustomerServiceMessage) {

	var dept struct {
		Department string                          `json:"department" bson:"department"`
		Queue      []models.CustomerServiceMessage `json:"queue" bson:"queue"`
	}

	details.ID = primitive.NewObjectID()
	deptName := "Other"
	filter := bson.M{"department": deptName}
	update := bson.M{"$push": bson.M{"queue": details}}

	cs := Database.Collection(os.Getenv("CUSTOMER_SERVICE_COLLECTION"))
	result := cs.FindOneAndUpdate(context.Background(), filter, update)

	if result.Err() == mongo.ErrNoDocuments {
		// DOCUMENT NOT FOUND
		log.Printf("Error grabbing dept: %s, Creating department now.... \n", dept.Department)
		dept.Department = deptName
		dept.Queue = append(dept.Queue, details)
		fmt.Print(dept)
		_, err := cs.InsertOne(context.Background(), dept)

		if err != nil {
			log.Printf("Error saving to dept: %s\n %v \n", deptName, err)
		}
	}
}

func SendToAccountDept(details models.CustomerServiceMessage) {

	var dept struct {
		Department string                          `json:"department" bson:"department"`
		Queue      []models.CustomerServiceMessage `json:"queue" bson:"queue"`
	}

	details.ID = primitive.NewObjectID()
	deptName := "Accounts"
	filter := bson.M{"department": deptName}
	update := bson.M{"$push": bson.M{"queue": details}}

	cs := Database.Collection(os.Getenv("CUSTOMER_SERVICE_COLLECTION"))
	result := cs.FindOneAndUpdate(context.Background(), filter, update)

	if result.Err() == mongo.ErrNoDocuments {
		// DOCUMENT NOT FOUND
		log.Printf("Error grabbing dept: %s, Creating department now.... \n", dept.Department)
		dept.Department = deptName
		dept.Queue = append(dept.Queue, details)
		fmt.Print(dept)
		_, err := cs.InsertOne(context.Background(), dept)

		if err != nil {
			log.Printf("Error saving to dept: %s\n %v \n", deptName, err)
		}
	}
}

func SendToBillingDept(details models.CustomerServiceMessage) {

	var dept struct {
		Department string                          `json:"department" bson:"department"`
		Queue      []models.CustomerServiceMessage `json:"queue" bson:"queue"`
	}
	deptName := "Billing"
	filter := bson.M{"department": deptName}
	update := bson.M{"$push": bson.M{"queue": details}}

	cs := Database.Collection(os.Getenv("CUSTOMER_SERVICE_COLLECTION"))
	result := cs.FindOneAndUpdate(context.Background(), filter, update)

	if result.Err() == mongo.ErrNoDocuments {
		// DOCUMENT NOT FOUND
		log.Printf("Error grabbing dept: %s, Creating department now.... \n", dept.Department)
		dept.Department = deptName
		dept.Queue = append(dept.Queue, details)
		fmt.Print(dept)
		_, err := cs.InsertOne(context.Background(), dept)

		if err != nil {
			log.Printf("Error saving to dept: %s\n %v \n", deptName, err)
		}
	}
}
