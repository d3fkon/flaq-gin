package payments

import (
	"net/http"
	"time"

	"github.com/d3fkon/gin-flaq/models"
	"github.com/d3fkon/gin-flaq/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Register a new payment for the user
func RegisterPayment(user models.User, amount string) interface{} {
	payment := models.Payment{
		User:      user.Id,
		Amount:    amount,
		Id:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
	}
	models.PaymentModel.New(payment)
	return payment
}

func GetAllPayments(user models.User) interface{} {
	query := bson.M{
		"User": user.Id,
	}
	elem := []models.Payment{}
	err := models.PaymentModel.FindMany(query, &elem)
	if err != nil {
		utils.Panic(http.StatusNotFound, "Cannot find payments", nil)
	}
	return elem
}
