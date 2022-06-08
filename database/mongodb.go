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

func InsertUserAndNotification(transaction models.Transaction) bool {
	fmt.Printf("Transaction: %v \n", transaction)
	return true
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