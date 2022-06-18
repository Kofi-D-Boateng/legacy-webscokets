package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"github.com/Kofi-D-Boateng/legacynotifications/utils"
)

func EmailHandler(w http.ResponseWriter, r *http.Request) {
	var variables models.EmailAttributes

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&variables)

	if err != nil {
		fmt.Print(err)
	}


	var result int = utils.SendConfirmationEmail(variables)
	w.WriteHeader(result)
	json.NewEncoder(w)
}