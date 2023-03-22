package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
)

func InsertToDatabase(payload json.RawMessage) (models.Response,error) {
	var notification models.TransactionNotificationVariables
	err := json.Unmarshal(payload, &notification)
	if err != nil {
		log.Println(err)
		return models.Response{StatusCode: http.StatusInternalServerError,Body:[]byte("")},err
	} else {
		err := InsertUserAndNotification(notification)
		if err != nil {
			return models.Response{StatusCode: http.StatusInternalServerError, Body: []byte("")},err
		}
		return models.Response{StatusCode: http.StatusOK,Body: []byte("")},nil
	}
}
