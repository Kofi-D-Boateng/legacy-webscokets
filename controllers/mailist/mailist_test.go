package controllers

import (
	"testing"

	"github.com/Kofi-D-Boateng/legacynotifications/utils"
)

func TestMailingListHandler(t *testing.T) {

	var data string = "kdboat2@gmail.com"
	var result int = utils.SendMailingListConfirmation(data)

	if result == 200 {
		t.Log("Passed test")
	} else {
		t.Error("Failed test")
	}

}