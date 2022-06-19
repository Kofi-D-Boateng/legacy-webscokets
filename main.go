package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Kofi-D-Boateng/legacynotifications/router"
	"github.com/joho/godotenv"
)

func main(){
	// GRABBING SET UP ENV VAR
	env := os.Getenv("GO_ENV")
	if env == "dev"{
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error: %s \n", err)
		}
	}
	
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