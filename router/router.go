package router

import (
	"fmt"
	"os"

	email "github.com/Kofi-D-Boateng/legacynotifications/controllers/email"
	mailist "github.com/Kofi-D-Boateng/legacynotifications/controllers/mailist"
	service "github.com/Kofi-D-Boateng/legacynotifications/controllers/service"
	user "github.com/Kofi-D-Boateng/legacynotifications/controllers/user"

	"github.com/gorilla/mux"
)


func Router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	api_version := os.Getenv("API_VERSION")

	userUri := fmt.Sprintf("/%s/user",api_version)
	setNotificationsUri := fmt.Sprintf("/%s/user/set-notifications",api_version)
	markNotificationsUri := fmt.Sprintf("/%s/user/mark-notifications",api_version)
	sendEmailUri := fmt.Sprintf("/%s/verification/send-email",api_version)
	customerServiceUri := fmt.Sprintf("/%s/customer-service/email-customer-service",api_version)
	mailListUri := fmt.Sprintf("/%s/mail-list/add",api_version)

	router.HandleFunc(userUri, user.GetNotificationsHandler).Queries("email", "{email}").Methods("GET")
	router.HandleFunc(setNotificationsUri, user.SetNotificationsHandler).Methods("PUT")
	router.HandleFunc(markNotificationsUri, user.MarkNotificationsHandler).Methods("PUT")
	// VERIFICATION
	router.HandleFunc(sendEmailUri, email.EmailHandler).Methods("POST")
	// CUSTOMER SERVICE
	router.HandleFunc(customerServiceUri, service.CustomerServiceHandler).Methods("PUT")
	// MAIL LIST SERVICE
	router.HandleFunc(mailListUri, mailist.MailingListHandler).Methods("PUT")


	return router
}