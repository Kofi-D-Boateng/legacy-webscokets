package database

import (
	"fmt"
	"testing"

	"github.com/Kofi-D-Boateng/legacynotifications/models"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockedFunc struct {
	mock.Mock
}

func (m* MockedFunc) FindAUser(email string) models.User{
	args := m.Called(email)

	fmt.Printf("ARGS: %v", args)
	user := models.User{
		ID: primitive.NewObjectID(),
		Email: args.String(),
		Notifications: []models.Transaction{},
	}
	return user
}

func TestFindUser(t *testing.T) {
	email := "kdboat2@gmail.com"
	MockedFindRepo := new(MockedFunc)
	MockedFindRepo.On("FindAUser", email).Return(models.User {
		ID: primitive.NewObjectID(),
		Email: email,
		Notifications: []models.Transaction{},
	})
	MockedFindRepo.FindAUser(email)
}