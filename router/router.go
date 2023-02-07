package router

import (
	"fmt"
	"os"

	"github.com/Kofi-D-Boateng/legacynotifications/controllers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	api_version := os.Getenv("API_VERSION")

	userUri := fmt.Sprintf("/%s/user", api_version)

	router.HandleFunc(userUri, controllers.GetNotificationsHandler).Queries("email", "{email}").Methods("GET")

	return router
}
