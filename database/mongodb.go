package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Db *mongo.Database
	UserCollection string
	CustomerServiceCollection string
)



func init(){
	_,file,_, ok := runtime.Caller(0)
	basePath := filepath.Dir(file)
	fmt.Println(file)
	fmt.Println(basePath)

	if !ok {
		log.Fatalf("Unable to find file path: %v", file)
	}

	err := godotenv.Load(filepath.Join(basePath, "../.env"))
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

