package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Kofi-D-Boateng/legacynotifications/router"
	"github.com/joho/godotenv"
)

func main(){
	// GRABBING SET UP ENV VAR
	err := godotenv.Load(".env")
	r := router.Router()
	if err != nil {
		log.Fatalf("Error: %s \n", err)
	}
	port := os.Getenv("PORT")
	
	
	fmt.Printf("Server listening at port%v \n", port)
	log.Fatal(http.ListenAndServe(port, r))
}