package database

import (
	"context"
	"log"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"go.mongodb.org/mongo-driver/bson"
)


func SendToOther(details models.CustomerServiceMessage) int16 {
	
	var dept struct {
		Department string `json:"department" bson:"department"`
		Queue		[]models.CustomerServiceMessage `json:"queue" bson:"queue"`
	}
	deptName := "Other"
	filter := bson.M{"department": deptName}

	cs := Db.Collection(CustomerServiceCollection)

	err := cs.FindOne(context.Background(), filter).Decode(&dept)
	
	if err != nil {
		log.Fatalf("Error grabbing collection: %s \n %v", CustomerServiceCollection, err)
		return 500
	}

	if dept.Department == deptName {
		dept.Queue = append(dept.Queue, details)
	}
	_, updateErr := cs.UpdateOne(context.Background(), filter, dept)

	if updateErr != nil {
		log.Fatalf("Error saving to dept: %s\n %v", deptName, updateErr)
		return 500
	}

	return 200

}

func SendToAccountDept(details models.CustomerServiceMessage) int16 {
	var dept struct {
		Department string `json:"department" bson:"department"`
		Queue		[]models.CustomerServiceMessage `json:"queue" bson:"queue"`
	}
	deptName := "Billing"
	filter := bson.M{"department": deptName}

	cs := Db.Collection(CustomerServiceCollection)

	err := cs.FindOne(context.Background(), filter).Decode(&dept)
	
	if err != nil {
		log.Fatalf("Error grabbing collection: %s \n %v", CustomerServiceCollection, err)
		return 500
	}

	if dept.Department == deptName {
		dept.Queue = append(dept.Queue, details)
	}
	_, updateErr := cs.UpdateOne(context.Background(), filter, dept)

	if updateErr != nil {
		log.Fatalf("Error saving to dept: %s\n %v", deptName, updateErr)
		return 500
	}



	return 200
}

func SendToBillingDept(details models.CustomerServiceMessage) int16{

	var dept struct {
		Department string `json:"department" bson:"department"`
		Queue		[]models.CustomerServiceMessage `json:"queue" bson:"queue"`
	}
	deptName := "Other"
	filter := bson.M{"department": deptName}

	cs := Db.Collection(CustomerServiceCollection)

	err := cs.FindOne(context.Background(), filter).Decode(&dept)
	
	if err != nil {
		log.Fatalf("Error grabbing collection: %s \n %v", CustomerServiceCollection, err)
		return 500
	}
	
	if dept.Department == deptName {
		dept.Queue = append(dept.Queue, details)
	}
	_, updateErr := cs.UpdateOne(context.Background(), filter, dept)

	if updateErr != nil {
		log.Fatalf("Error saving to dept: %s\n %v", deptName, updateErr)
		return 500
	}

	return 200
}