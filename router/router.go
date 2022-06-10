package router

import (
	emailer "github.com/Kofi-D-Boateng/legacynotifications/controllers/emailer"
	mailist "github.com/Kofi-D-Boateng/legacynotifications/controllers/mailist"
	service "github.com/Kofi-D-Boateng/legacynotifications/controllers/service"
	user "github.com/Kofi-D-Boateng/legacynotifications/controllers/user"

	"github.com/gorilla/mux"
)


func Router() *mux.Router {
	router := mux.NewRouter()
	// NOTIFICATIONS
	router.HandleFunc("/api/v1/user", user.GetNotificationsHandler).Methods("GET")
	router.HandleFunc("/api/v1/user/set-notifications", user.SetNotificationsHandler).Methods("PUT")
	router.HandleFunc("/api/v1/user/mark-notification", user.MarkNotificationsHandler).Methods("PUT")
	// VERIFICATION
	router.HandleFunc("/api/v1/verification/send-email", emailer.EmailHandler).Methods("POST")
	router.HandleFunc("/api/v1/verification/new-verification-link", emailer.EmailHandler).Methods("POST")
	// CUSTOMER SERVICE
	router.HandleFunc("/api/v1/customer-service/email-customer-service", service.CustomerServiceHandler).Methods("PUT")
	// MAIL LIST SERVICE
	router.HandleFunc("/api/v1/mail-list/add", mailist.MailingListHandler).Methods("PUT")


	return router
}