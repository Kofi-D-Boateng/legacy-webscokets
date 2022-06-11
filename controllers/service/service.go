package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/Kofi-D-Boateng/legacynotifications/database"
	"github.com/Kofi-D-Boateng/legacynotifications/models"
) 

var EmailExpr *regexp.Regexp
var Accounts *regexp.Regexp
var Billing *regexp.Regexp

func init(){
	email, checkOne := regexp.Compile(`.+@.+\..+`)
	account, checkTwo := regexp.Compile(`(?i)\\Account.$|Accounts.$|Transfer.$\\`)
	billing, checkThree := regexp.Compile(`(?i)\\billing.$|notice.$|\\`)

	if checkOne != nil  {
		log.Fatalf("Error with compiling regex partner for %s:  %v \n",email, checkOne)
	}
	if checkTwo != nil  {
		log.Fatalf("Error with compiling regex partner for %s:  %v \n",account, checkTwo)
	}
	if checkThree != nil  {
		log.Fatalf("Error with compiling regex partner for %s:  %v \n",billing, checkThree)
	}


	EmailExpr = email
	Accounts = account
	Billing = billing
}


func CustomerServiceHandler(w http.ResponseWriter, r *http.Request){
	var customerService models.CustomerServiceMessage
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&customerService)
	if err != nil {
		log.Fatal(err)
	}

	if !EmailExpr.Match([]byte(customerService.Emailer)) {
		log.Fatalf("Error with email: %s", customerService.Emailer)
	}

	if Accounts.Match([]byte(customerService.Topic)){
		var result int = database.SendToAccountDept(customerService)
		w.WriteHeader(result)
		json.NewEncoder(w)
	}

	if Billing.Match([]byte(customerService.Topic)){
		var result int = database.SendToBillingDept(customerService)
		w.WriteHeader(result)
		json.NewEncoder(w)
	}

	var result int = database.SendToOther(customerService)
	w.WriteHeader(result)
	json.NewEncoder(w)
}