package controllers

import (
	"testing"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"github.com/Kofi-D-Boateng/legacynotifications/utils"
)

func TestEmailHandler(t *testing.T){

	var data models.EmailAttributes
	data.Person.Email = "kdboat2@gmail.com"
	data.Person.Name = "Kofi Boateng"
	data.Token = "12das1F!2gdsf#1"

	var result int = utils.SendConfirmationEmail(data)

	if result == 200 {
		t.Log("Passed test")
	}else{
		t.Error("Failed test")
	}

}
