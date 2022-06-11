package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Kofi-D-Boateng/legacynotifications/router"
)

func main(){
	// GRABBING SET UP ENV VAR
	// err := godotenv.Load(".env")
	r := router.Router()
	// if err != nil {
	// 	log.Fatalf("Error: %s", err)
	// }
	port := os.Getenv("PORT")
	
	
	fmt.Printf("Server listening at port%v", port)
	log.Fatal(http.ListenAndServe(port, r))
}