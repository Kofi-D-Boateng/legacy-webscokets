package database

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


func SendToOther(details models.CustomerServiceMessage) int {
	
	var dept struct {
		Department string `json:"department" bson:"department"`
		Queue		[]models.CustomerServiceMessage `json:"queue" bson:"queue"`
	}

	details.ID = primitive.NewObjectID()
	deptName := "Other"
	filter := bson.M{"department": deptName}
	update := bson.M{"$push": bson.M{"queue":details}}
	
	cs := Database.Db.Collection(Database.CustomerServiceCollection)
	result := cs.FindOneAndUpdate(context.Background(), filter, update)

	if result.Err() == mongo.ErrNoDocuments {
		// DOCUMENT NOT FOUND
		log.Printf("Error grabbing dept: %s, Creating department now.... \n", dept.Department)
		dept.Department = deptName
		dept.Queue = append(dept.Queue, details)
		fmt.Print(dept)
		_, err := cs.InsertOne(context.Background(),dept)
		
		if err != nil {
			log.Printf("Error saving to dept: %s\n %v \n", deptName, err)
			return http.StatusInternalServerError
		}
		
		return http.StatusOK
	}
	
	return http.StatusOK

}

func SendToAccountDept(details models.CustomerServiceMessage) int {


	var dept struct {
		Department string `json:"department" bson:"department"`
		Queue		[]models.CustomerServiceMessage `json:"queue" bson:"queue"`
	}

	details.ID = primitive.NewObjectID()
	deptName := "Accounts"
	filter := bson.M{"department": deptName}
	update := bson.M{"$push": bson.M{"queue":details}}
	
	cs := Database.Db.Collection(Database.CustomerServiceCollection)
	result := cs.FindOneAndUpdate(context.Background(), filter, update)

	
	if result.Err() == mongo.ErrNoDocuments {
		// DOCUMENT NOT FOUND
		log.Printf("Error grabbing dept: %s, Creating department now.... \n", dept.Department)
		dept.Department = deptName
		dept.Queue = append(dept.Queue, details)
		fmt.Print(dept)
		_, err := cs.InsertOne(context.Background(),dept)
		
		if err != nil {
			log.Printf("Error saving to dept: %s\n %v \n", deptName, err)
			return http.StatusInternalServerError
		}
		
		return http.StatusOK
	}
	
	return http.StatusOK
}

func SendToBillingDept(details models.CustomerServiceMessage) int{

	var dept struct {
		Department string `json:"department" bson:"department"`
		Queue		[]models.CustomerServiceMessage `json:"queue" bson:"queue"`
	}
	deptName := "Billing"
	filter := bson.M{"department": deptName}
	update := bson.M{"$push": bson.M{"queue":details}}
	
	cs := Database.Db.Collection(Database.CustomerServiceCollection)
	result := cs.FindOneAndUpdate(context.Background(), filter, update)

	
	if result.Err() == mongo.ErrNoDocuments {
		// DOCUMENT NOT FOUND
		log.Printf("Error grabbing dept: %s, Creating department now.... \n", dept.Department)
		dept.Department = deptName
		dept.Queue = append(dept.Queue, details)
		fmt.Print(dept)
		_, err := cs.InsertOne(context.Background(),dept)
		
		if err != nil {
			log.Printf("Error saving to dept: %s\n %v \n", deptName, err)
			return http.StatusInternalServerError
		}
		
		return http.StatusOK
	}
	
	return http.StatusOK
}