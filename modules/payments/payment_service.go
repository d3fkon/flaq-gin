package payments

import (
	"net/http"
	"strconv"
	"time"

	"github.com/d3fkon/gin-flaq/models"
	"github.com/d3fkon/gin-flaq/modules/users"
	"github.com/d3fkon/gin-flaq/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Register a new payment for the user
func RegisterPayment(user models.User, amount string) interface{} {
	amountInt, _ := strconv.Atoi(amount)
	flaqReward := getFlaqForPayment(amountInt)
	flaqRewardStr := strconv.Itoa(flaqReward)
	payment := models.Payment{
		User:       user.Id,
		Amount:     amount,
		Id:         primitive.NewObjectID(),
		CreatedAt:  primitive.NewDateTimeFromTime(time.Now()),
		FlaqReward: flaqRewardStr,
	}
	models.PaymentModel.New(payment)
	// Update user's flaq balance
	users.UpdateFlaqPoints(&user, flaqReward)
	return payment
}

func getFlaqForPayment(amount int) int {
	if amount > 100 {
		return 100
	} else {
		return amount
	}
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
