package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Database
var UserCollection string;
var CustomerServiceCollection string;



func init(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// GRABBING DB INFO
	var dns string = os.Getenv("MONGO_URI")
	var dbName string = os.Getenv("DB_NAME")
	UserCollection = os.Getenv("USERS_COLLECTION")
	CustomerServiceCollection = os.Getenv("CUSTOMER_SERVICE_COLLECTION")

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

