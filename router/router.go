package router

import (
	"github.com/Kofi-D-Boateng/legacynotifications/controllers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	// NOTIFICATIONS
	router.HandleFunc("/api/v1/user", controllers.GetNotificationsHandler).Methods("GET")
	router.HandleFunc("/api/v1/user/set-notifications", controllers.SetNotificationsHandler).Methods("PUT")
	router.HandleFunc("/api/v1/user/mark-notification", controllers.MarkNotificationsHandler).Methods("PUT")
	// VERIFICATION
	router.HandleFunc("/api/v1/verification/send-email", controllers.EmailHandler).Methods("POST")
	router.HandleFunc("/api/v1/verification/new-verification-link", controllers.EmailHandler).Methods("POST")
	// CUSTOMER SERVICE
	router.HandleFunc("/api/v1/customer-service/email-customer-service", controllers.CustomerServiceHandler).Methods("PUT")
	// MAIL LIST SERVICE
	router.HandleFunc("/api/v1/mail-list/add", controllers.MailingListHandler).Methods("PUT")


	return router
}