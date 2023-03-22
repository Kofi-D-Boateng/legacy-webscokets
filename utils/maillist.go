package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
)

func AddToMailList(payload json.RawMessage) (models.Response,error) {
	var maillist struct {
		Email string `json:"email"`
	}

	err := json.Unmarshal(payload, &maillist.Email)
	if err != nil {
		log.Printf("Error unmarshalling message: %s", err)
		return models.Response{StatusCode: http.StatusInternalServerError,Body: []byte("")},err
	} else {
		fmt.Printf("EMAIL: %v\n", maillist.Email)
		err := SendMailingListConfirmation(maillist.Email)
		if err != nil{
			return models.Response{StatusCode: http.StatusInternalServerError,Body: []byte("")},err
		}
	}
	return models.Response{StatusCode: http.StatusOK,Body: []byte("")},nil
}
