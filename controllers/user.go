package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/utils"
)

func GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	email := r.FormValue("email")
	foundUser := utils.FindAUser(email)

	json.NewEncoder(w).Encode(foundUser.Notifications)
}
