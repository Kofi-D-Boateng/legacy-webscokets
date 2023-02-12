package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"github.com/Kofi-D-Boateng/legacynotifications/router"
	"github.com/Kofi-D-Boateng/legacynotifications/utils"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func init() {
	os.Setenv("GO_ENV", "dev")
	env := os.Getenv("GO_ENV")
	if env == "dev" {
		_, file, _, ok := runtime.Caller(0)
		basePath := filepath.Dir(file)
		fmt.Println(file)
		fmt.Println(basePath)

		if !ok {
			log.Fatalf("Unable to find file path: %v", file)
		}

		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	}

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

	utils.EmailExpr = email
	utils.Accounts = account
	utils.Billing = billing
}

func main() {

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_CONN"))
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	utils.StartMaillistQueue(conn)
	utils.StartUpdateQueue(conn)
	utils.StartInsertQueue(conn)
	utils.StartVerificationQueue(conn)
	utils.StartCustomerServiceQueue(conn)

	r := router.Router()
	port := os.Getenv("PORT")
	utils.ConnectDatabase(os.Getenv("MONGO_URI"), os.Getenv("DB_NAME"))

	srv := &http.Server{
		Handler:      r,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Server listening at port%v \n", port)

	log.Fatal(srv.ListenAndServe())
}
