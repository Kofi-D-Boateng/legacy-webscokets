package database

import (
	"context"
	"log"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"go.mongodb.org/mongo-driver/bson"
)


func SendToOther(details models.CustomerServiceMessage) int {
	
	var dept struct {
		Department string `json:"department" bson:"department"`
		Queue		[]models.CustomerServiceMessage `json:"queue" bson:"queue"`
	}
	deptName := "Other"
	filter := bson.M{"department": deptName}

	cs := Db.Collection(CustomerServiceCollection)

	err := cs.FindOne(context.Background(), filter).Decode(&dept)
	
	if err != nil {
		log.Printf("Error grabbing collection: %s \n %v \n", CustomerServiceCollection, err)
		return http.StatusInternalServerError
	}

	if dept.Department == deptName {
		dept.Queue = append(dept.Queue, details)
	}
	_, updateErr := cs.UpdateOne(context.Background(), filter, dept)

	if updateErr != nil {
		log.Printf("Error saving to dept: %s\n %v \n", deptName, updateErr)
		return http.StatusInternalServerError
	}

	return http.StatusOK

}

func SendToAccountDept(details models.CustomerServiceMessage) int {
	var dept struct {
		Department string `json:"department" bson:"department"`
		Queue		[]models.CustomerServiceMessage `json:"queue" bson:"queue"`
	}
	deptName := "Billing"
	filter := bson.M{"department": deptName}

	cs := Db.Collection(CustomerServiceCollection)

	err := cs.FindOne(context.Background(), filter).Decode(&dept)
	
	if err != nil {
		log.Printf("Error grabbing collection: %s \n %v", CustomerServiceCollection, err)
		return http.StatusInternalServerError
	}

	if dept.Department == deptName {
		dept.Queue = append(dept.Queue, details)
	}
	_, updateErr := cs.UpdateOne(context.Background(), filter, dept)

	if updateErr != nil {
		log.Printf("Error saving to dept: %s\n %v", deptName, updateErr)
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func SendToBillingDept(details models.CustomerServiceMessage) int{

	var dept struct {
		Department string `json:"department" bson:"department"`
		Queue		[]models.CustomerServiceMessage `json:"queue" bson:"queue"`
	}
	deptName := "Other"
	filter := bson.M{"department": deptName}

	cs := Db.Collection(CustomerServiceCollection)

	err := cs.FindOne(context.Background(), filter).Decode(&dept)
	
	if err != nil {
		log.Printf("Error grabbing collection: %s \n %v", CustomerServiceCollection, err)
		return http.StatusInternalServerError
	}
	
	if dept.Department == deptName {
		dept.Queue = append(dept.Queue, details)
	}
	_, updateErr := cs.UpdateOne(context.Background(), filter, dept)

	if updateErr != nil {
		log.Printf("Error saving to dept: %s\n %v \n", deptName, updateErr)
		return http.StatusInternalServerError
	}

	return http.StatusOK
}