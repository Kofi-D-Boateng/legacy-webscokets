package models

import (
	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Decoder = schema.NewDecoder()


type User struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email         string `json:"email,omitempty" bson:"email,omitempty"`
	Notifications []Transaction `json:"notifications,omitempty" bson:"notifications,omitempty"`
}

type Transaction struct {
	Receiver				string `json:"receiver" bson:"cardNumber"`
	ReceiverEmail			string `json:"receiverEmail" bson:"cardNumber"`
	Sender					string `json:"sender" bson:"cvc"`
	DateOfTransaction 		string `json:"dateOfTransaction" bson:"dateOfTransaction"`
	Type 					string `json:"type" bson:"type"`
	Amount 					float64 `json:"amount" bson:"amount"`
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