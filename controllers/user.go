package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"github.com/Kofi-D-Boateng/legacynotifications/utils"
)

func GetNotificationsHandler(payload json.RawMessage)(models.Response,error) {
	
	var email string
	err := json.Unmarshal(payload,&email)
	if err != nil{
		return models.Response{StatusCode: http.StatusBadRequest,Body: []byte("")},err
	}
	
	foundUser,err := utils.FindAUser(email)
	if err != nil{
		return models.Response{StatusCode: http.StatusUnauthorized,Body: []byte("")},err
	}
	user,_ := json.Marshal(foundUser)
	return models.Response{StatusCode: http.StatusOK,Body: user},err
}
