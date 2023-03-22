package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
)

func SendVerificationEmail(payload json.RawMessage) (models.Response,error) {

	var emailAttributes models.EmailAttributes
	err := json.Unmarshal(payload, &emailAttributes)
	if err != nil {
		log.Printf("Error unmarshalling message: %s", err)
		return models.Response{StatusCode: http.StatusInternalServerError,Body: []byte("")},err
	}
	er := SendConfirmationEmail(emailAttributes)
	if er != nil {
		return models.Response{StatusCode: http.StatusInternalServerError,Body: []byte("")},err
	}
	return models.Response{StatusCode: http.StatusOK,Body: []byte("")},nil
}
