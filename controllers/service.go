package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
)



func CustomerServiceHandler(w http.ResponseWriter, r *http.Request){
	var customerService models.CustomerServiceMessage
	err := models.Decoder.Decode(&customerService,r.URL.Query())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n json: %v", customerService)
}