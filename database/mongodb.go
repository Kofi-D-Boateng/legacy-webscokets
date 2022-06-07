package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Database

func init(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	var dns string = os.Getenv("MONGO_URI")
	var dbName string = os.Getenv("DB_NAME")
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

func Insert(user models.User) bool {
	fmt.Printf("got email: %v \n", user)
	return true
}

func Find(email string) models.User{
	fmt.Printf("\n got email: %v \n", email)
	return models.User{}
}

func Update() bool {
	return true
}