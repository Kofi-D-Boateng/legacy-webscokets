package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/utils"
)

func MailingListHandler(w http.ResponseWriter, r *http.Request) {
	var email string;
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&email)

	if err != nil {
		log.Fatal(err)
	}

	var result int =  utils.SendMailingListConfirmation(email)
	w.WriteHeader(result)
	json.NewEncoder(w)
}