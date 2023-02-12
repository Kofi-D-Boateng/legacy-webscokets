package utils

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"github.com/joho/godotenv"
)

var (
	companyEmail string
	link         string
	password     string
	smptHost     string
	smtpPort     string
)

func init() {
	env := os.Getenv("GO_ENV")
	if env == "dev" {

		_, file, _, ok := runtime.Caller(0)
		basePath := filepath.Dir(file)
		fmt.Println(file)
		fmt.Println(basePath)

		if !ok {
			log.Fatalf("Unable to find file path: %v", file)
		}

		err := godotenv.Load(filepath.Join(basePath, "../.env"))
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	}

	var from string = os.Getenv("COMPANY_EMAIL")
	var pw string = os.Getenv("COMPANY_PASSWORD")
	var host string = os.Getenv("SMTP_HOST")
	var port string = os.Getenv("SMTP_PORT")
	var accountAuthLink string = os.Getenv("ACCT_AUTH_LINK")

	password = pw
	companyEmail = from
	link = accountAuthLink
	smptHost = host
	smtpPort = port

}

type loginAuthStruct struct {
	username, password string
}

func loginAuth(username, password string) smtp.Auth {
	return &loginAuthStruct{username, password}
}

func (a *loginAuthStruct) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuthStruct) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unkown fromserver")
		}
	}
	return nil, nil
}

func SendConfirmationEmail(attributes models.EmailAttributes) {

	auth := loginAuth(companyEmail, password)

	var msg string
	if attributes.Person.Name != "" {
		msg = fmt.Sprintf("Thank you %s for opening an account an account with Legacy Bank. Please click on the link below to verify your account. %s?token=%s", attributes.Person.Name, link, attributes.Token)
	} else {
		msg = fmt.Sprintf("You have requested a new link to verify your account. Please click on the link below to verify your account. %s?token=%s", link, attributes.Token)
	}

	from := fmt.Sprintf("From: <%s>\r\n", companyEmail)
	to := fmt.Sprintf("To: <%s>\r\n", attributes.Person.Name)
	subject := "Thank you for joining the mailing list\r\n"
	body := msg + "\r\n"

	message := from + to + subject + "\r\n" + body

	err := smtp.SendMail(smptHost+":"+smtpPort, auth, companyEmail, []string{attributes.Person.Email}, []byte(message))

	if err != nil {
		fmt.Printf("Cannot send email. Error: %v", err)
	}
}

func SendMailingListConfirmation(email string) {
	auth := loginAuth(companyEmail, password)

	from := fmt.Sprintf("From: <%s>\r\n", companyEmail)
	to := fmt.Sprintf("To: <%s>\r\n", email)
	subject := "Thank you for joining the mailing list\r\n"
	body := "Thank you for joining submitting your email to the maillist! This is a fake end point for demonstration therefore your email will not be stored.\r\n-Kofi Boateng\r\n"

	msg := from + to + subject + "\r\n" + body

	err := smtp.SendMail(smptHost+":"+smtpPort, auth, companyEmail, []string{email}, []byte(msg))

	if err != nil {
		log.Println(err)
	}
}
