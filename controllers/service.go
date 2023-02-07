package controllers

import (
	"log"
	"regexp"

	"github.com/Kofi-D-Boateng/legacynotifications/database"
	"github.com/Kofi-D-Boateng/legacynotifications/models"
)

var (
	EmailExpr *regexp.Regexp
	Accounts  *regexp.Regexp
	Billing   *regexp.Regexp
)

func init() {
	email, checkOne := regexp.Compile(`.+@.+\..+`)
	account, checkTwo := regexp.Compile(`(?i)\\Account.$|Accounts.$|Transfer.$\\`)
	billing, checkThree := regexp.Compile(`(?i)\\billing.$|notice.$|\\`)

	if checkOne != nil {
		log.Fatalf("Error with compiling regex partner for %s:  %v \n", email, checkOne)
	}
	if checkTwo != nil {
		log.Fatalf("Error with compiling regex partner for %s:  %v \n", account, checkTwo)
	}
	if checkThree != nil {
		log.Fatalf("Error with compiling regex partner for %s:  %v \n", billing, checkThree)
	}

	EmailExpr = email
	Accounts = account
	Billing = billing
}

func CustomerServiceRequest(customerService models.CustomerServiceMessage) {

	if !EmailExpr.Match([]byte(customerService.Emailer)) {
		log.Fatalf("Error with email: %s", customerService.Emailer)
	}

	if Accounts.Match([]byte(customerService.Topic)) {
		database.SendToAccountDept(customerService)
	}

	if Billing.Match([]byte(customerService.Topic)) {
		database.SendToBillingDept(customerService)
	}

	database.SendToOther(customerService)
}
