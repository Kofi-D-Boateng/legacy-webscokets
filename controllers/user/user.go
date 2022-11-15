package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/database"
	"github.com/Kofi-D-Boateng/legacynotifications/models"
)


func GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	email := r.FormValue("email")	
	foundUser := database.FindAUser(email)
  
	json.NewEncoder(w).Encode(foundUser)
}

func SetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	var transactionNotificationRequest models.TransactionNotificationVariables
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&transactionNotificationRequest)
	fmt.Printf("\n query: %v\n", transactionNotificationRequest)

	if err != nil {
		log.Fatal(err)
	}


	var result int = database.InsertUserAndNotification(transactionNotificationRequest)
	w.WriteHeader(result)
	json.NewEncoder(w)
}

func MarkNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	var MessageToMark models.MarkMessage
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&MessageToMark)
	if err != nil {
		log.Fatal(err)
	}

	var result models.User = database.MarkMessageAsRead(MessageToMark)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}