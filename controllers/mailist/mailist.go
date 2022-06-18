package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/utils"
)

func MailingListHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&data)

	if err != nil {
		log.Fatal(err)
	}

	var result int =  utils.SendMailingListConfirmation(data.Email)
	w.WriteHeader(result)
	json.NewEncoder(w)
}