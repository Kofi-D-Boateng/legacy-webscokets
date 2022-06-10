package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/database"
	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"github.com/gorilla/schema"
)


func GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	decoder := schema.NewDecoder()
	err := decoder.Decode(&user, r.URL.Query())
	if err != nil {
		log.Fatal(err)
	}
	foundUser := database.FindAUser(user.Email)
	json.NewEncoder(w).Encode(foundUser)
}

func SetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var variables struct {
		Email					string 	`json:"email"`
		Receiver				string 	`json:"receiver" `
		ReceiverEmail			string 	`json:"receiverEmail"`
		Sender					string 	`json:"sender"`
		IsReceiverInDatabase 	bool 	`json:"isReceiverInDatabase"`
		DateOfTransaction 		string 	`json:"dateOfTransaction"`
		Type 					string 	`json:"type"`
		Amount 					float64 `json:"amount"`
	}
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&variables)
	fmt.Printf("\n query: %v\n", variables)

	if err != nil {
		log.Fatal(err)
	}

	var result bool = database.InsertUserAndNotification(variables)
	json.NewEncoder(w).Encode(result)
}

func MarkNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var variables models.MarkMessage
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&variables)
	if err != nil {
		log.Fatal(err)
	}

	var result bool = database.MarkMessageAsRead(variables)
	json.NewEncoder(w).Encode(result)
}