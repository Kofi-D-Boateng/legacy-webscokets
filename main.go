package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Kofi-D-Boateng/legacynotifications/database"
	"github.com/Kofi-D-Boateng/legacynotifications/router"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func init(){

	env := os.Getenv("GO_ENV")
	if env == "dev" {
		_,file,_, ok := runtime.Caller(0)
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
	client, err := mongo.Connect(context.TODO(),clientOptions)

	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("MongoDB connected to DB: %v \n", dbName)
	// Get DB
	database.Database.Db = client.Database(dbName)
}


func main(){
	// GRABBING SET UP ENV VAR
	
	r := router.Router()
	port := os.Getenv("PORT")
	
	srv := &http.Server{
		Handler: r,
		Addr: port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15* time.Second,
	}
	
	fmt.Printf("Server listening at port%v \n", port)
  
	log.Fatal(srv.ListenAndServe())
}