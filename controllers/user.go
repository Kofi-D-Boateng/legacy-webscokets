package controllers

import (
	"net/http"
	"os"

	"github.com/Kofi-D-Boateng/legacynotifications/database"
)

var usersCollection string = os.Getenv("USERS_COLLECTION")

func GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email := r.URL.Query().Get("email")
	go database.Find(email)




}

func SetNotificationsHandler(w http.ResponseWriter, r *http.Request) {

}

func MarkNotificationsHandler(w http.ResponseWriter, r *http.Request) {

}