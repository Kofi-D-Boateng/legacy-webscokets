package models

import (
	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Decoder = schema.NewDecoder()


type User struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email         string `json:"email,omitempty" bson:"email,omitempty"`
	Notifications []bson.M `json:"notifications,omitempty" bson:"notifications,omitempty"`
}

type Transaction struct {
	ID            		primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AccountNumber 		string `json:"accountNumber" bson:"accountNumber,omitempty"`
	CardNumber 			uint64 `json:"cardNumber" bson:"cardNumber,omitempty"`
	Cvc 				uint64 `json:"cvc" bson:"cvc,omitempty"`
	EmailOfTransferee 	string `json:"emailOfTransferee" bson:"emailOfTransferee,omitempty"`
	DateOfTransaction 	primitive.DateTime `json:"dateOfTransaction" bson:"dateOfTransaction,omitempty"`
	Type 				string `json:"type" bson:"type,omitempty"`
	Amount 				float64 `json:"amount" bson:"amount,omitempty"`
	PhoneNumber 		string `json:"phoneNumber" bson:"phoneNumber,omitempty"`
}

type MarkMessage struct {
	Email 		string `json:"email"`
	MsgID 		string `json:"msg_id"`
}

type CustomerServiceMessage struct {
	ID            		primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Emailer 			string `json:"emailer,omitempty" bson:"emailer,omitempty"`
	Topic				string `json:"topic,omitempty" bson:"_id,omitempty"`
	Message				string `json:"msg,omitempty" bson:"msg,omitempty"`
	SentAt				string `json:"sentAt,omitempty" bson:"sentAt,omitempty"` 
}