package router

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RequestResponse struct {
	name string
	server *httptest.Server
	response *models.User
}

func TestRequest(t *testing.T){

	tests := []RequestResponse {
		{
			name: "Basic Request",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"_id:62747bff214677526c771f6a , email: kdboat2@gmail.com, notifications:[{_id: 62747bff214677526c771f6b , receiver:john doe, sender:kofi boateng, date:1-04-2022, read:true, amount:120.54}]"}`))
			})),
			response: &models.User{
				ID: primitive.ObjectID{},
				Email: "kdboat2@gmail.com",
				Notifications: []models.Transaction{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T){
			defer test.server.Close()
			resp := http.Client{}
			if !reflect.DeepEqual(resp, test.response){
				t.Error("Did not work")
			}else{
				t.Log("Passed")
			}
		})
	}
}

func TestMain(t *testing.T){
	t.Setenv("API_VERSION", "api/v1")
	t.Setenv("PORT", ":5500")


	c := http.Client{}

	r,_ := c.Get("http://localhost:5500/api/v1/user?email=kdboat2@gmail.com")
	if r != nil {
		t.Logf("PASSED! Got: %v",r)
	}
}