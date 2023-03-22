package utils

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
)

var (
	CompanyEmail string
	Link         string
	Password     string
	SmptHost     string
	SmtpPort     string
)

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

func SendConfirmationEmail(attributes models.EmailAttributes) error {

	auth := loginAuth(CompanyEmail, Password)

	var msg string
	if attributes.Person.Name != "" {
		msg = fmt.Sprintf("Thank you %s for opening an account an account with Legacy Bank. Please click on the link below to verify your account. %s?token=%s", attributes.Person.Name, Link, attributes.Token)
	} else {
		msg = fmt.Sprintf("You have requested a new link to verify your account. Please click on the link below to verify your account. %s?token=%s", Link, attributes.Token)
	}

	from := fmt.Sprintf("From: <%s>\r\n", CompanyEmail)
	to := fmt.Sprintf("To: <%s>\r\n", attributes.Person.Name)
	subject := "Thank you for joining the mailing list\r\n"
	body := msg + "\r\n"

	message := from + to + subject + "\r\n" + body

	err := smtp.SendMail(SmptHost+":"+SmtpPort, auth, CompanyEmail, []string{attributes.Person.Email}, []byte(message))

	if err != nil {
		fmt.Printf("Cannot send email. Error: %v", err)
		return err
	}
	return nil
}

func SendMailingListConfirmation(email string) error {
	auth := loginAuth(CompanyEmail, Password)

	from := fmt.Sprintf("From: <%s>\r\n", CompanyEmail)
	to := fmt.Sprintf("To: <%s>\r\n", email)
	subject := "Thank you for joining the mailing list\r\n"
	body := "Thank you for joining submitting your email to the maillist! This is a fake end point for demonstration therefore your email will not be stored.\r\n-Kofi Boateng\r\n"

	msg := from + to + subject + "\r\n" + body

	err := smtp.SendMail(SmptHost+":"+SmtpPort, auth, CompanyEmail, []string{email}, []byte(msg))

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
