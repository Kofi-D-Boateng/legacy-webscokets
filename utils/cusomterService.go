package utils

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
)

func SendToCustomerService(payload json.RawMessage) error {

	var request models.CustomerServiceMessage
	err := json.Unmarshal(payload, &request)
	if err != nil {
		log.Printf("Error unmarshalling message: %s", err)
		return err;
	} else {
		if !EmailExpr.Match([]byte(request.Emailer)) {
			log.Fatalf("Error with email: %s", request.Emailer)
			return errors.New("email does not match predefined regex pattern")
		}
		if Accounts.Match([]byte(request.Topic)) {
			return SendToAccountDept(request)
		}
		if Billing.Match([]byte(request.Topic)) {
			return SendToBillingDept(request)
		}
		return SendToOther(request)
	}
}
