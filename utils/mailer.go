package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"github.com/joho/godotenv"
)


var host string
var port string
var companyEmail string
var companyPassword string
var link string

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	var from string = os.Getenv("COMPANY_EMAIL")
	var password string = os.Getenv("COMPANY_PASSWORD")
	var smtpHost string = os.Getenv("SMTP_HOST")
	var smtpPort string = os.Getenv("SMTP_HOST")
	var accountAuthLink string = os.Getenv("ACCT_AUTH_LINK")

	companyPassword = password
	companyEmail = from
	host = smtpHost
	port = smtpPort
	link = accountAuthLink
}

func SendConfirmationEmail(attributes models.EmailAttributes) int16 {

	to := []string{attributes.Person.Email}
	auth := smtp.PlainAuth("", companyEmail, companyPassword, host)

	var msg string 
	if attributes.Person.Name != "" {
		msg = fmt.Sprintf("Thank you %s for opening an account an account with Legacy Bank. Please click on the link below to verify your account. %s?token=%s",attributes.Person.Name, link, attributes.Token)
	}else {
		msg = fmt.Sprintf("You have requested a new link to verify your account. Please click on the link below to verify your account. %s?token=%s",link,attributes.Token)
	}


	err := smtp.SendMail(host+":"+port, auth,companyEmail,to, []byte(msg))

	if err != nil {
		log.Fatal(err)
		return 500
	}

	return 200

}
