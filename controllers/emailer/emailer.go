package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"github.com/Kofi-D-Boateng/legacynotifications/utils"
)

func EmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var variables models.EmailAttributes

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&variables)

	if err != nil {
		log.Fatal(err)
	}


	var result int16 = utils.SendConfirmationEmail(variables)
	json.NewEncoder(w).Encode(result)
}