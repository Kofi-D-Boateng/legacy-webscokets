package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/Kofi-D-Boateng/legacynotifications/controllers"
	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"github.com/Kofi-D-Boateng/legacynotifications/utils"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
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


	var from string = os.Getenv("COMPANY_EMAIL")
	var pw string = os.Getenv("COMPANY_PASSWORD")
	var host string = os.Getenv("SMTP_HOST")
	var port string = os.Getenv("SMTP_PORT")
	var accountAuthLink string = os.Getenv("ACCT_AUTH_LINK")

	utils.Password = pw
	utils.CompanyEmail = from
	utils.Link = accountAuthLink
	utils.SmptHost = host
	utils.SmtpPort = port


}

func main() {
	utils.ConnectDatabase(os.Getenv("MONGO_URI"), os.Getenv("DB_NAME"))
	lambda.Start(handler)
}

func handler(ctx context.Context, request models.Request) (models.Response,error){
	fmt.Printf("Request --> %v\n",request)
		
		switch request.Function{
			case "getNotifications":
				return controllers.GetNotificationsHandler(request.Payload)
			case "customerService":
				err := utils.SendToCustomerService(request.Payload)
				if err != nil{
					return models.Response{StatusCode: http.StatusInternalServerError},err
				}
				return models.Response{StatusCode: http.StatusOK,Body: []byte("")},nil
			case "insertNotifications":
				return utils.InsertToDatabase(request.Payload)
			case "addToMailList":
				return utils.AddToMailList(request.Payload)	
			case "updateNotification":
				return utils.UpdateNotifications(request.Payload)
			case "verificationEmail":
				return utils.SendVerificationEmail(request.Payload)
			default:
				return models.Response{},errors.New("unknown function")				
		}
}