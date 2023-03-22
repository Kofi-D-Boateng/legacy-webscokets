package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
)

func UpdateNotifications(payload json.RawMessage) (models.Response,error) {

	var markMessage models.MarkMessage
	err := json.Unmarshal(payload, &markMessage)
	if err != nil {
		log.Printf("Error unmarshalling message: %s", err)
		return models.Response{StatusCode: http.StatusInternalServerError,Body: []byte("")},err
	}
	user,err := MarkMessageAsRead(markMessage)
	if err != nil {
		return models.Response{StatusCode: http.StatusInternalServerError,Body:[]byte("")},err
	}
	u,_:= json.Marshal(user)
	return models.Response{StatusCode: http.StatusOK,Body: u},nil
}
