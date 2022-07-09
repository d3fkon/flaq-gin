package models

import (
	"github.com/d3fkon/gin-flaq/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	Id         primitive.ObjectID `bson:"_id"`
	User       primitive.ObjectID `bson:"User" binding:"required"`
	Amount     string             `bson:"Amount" binding:"required"`
	CreatedAt  primitive.DateTime `bson:"CreatedAt"`
	FlaqReward float64            `bson:"FlaqReward"`
}

var PaymentModel = Collection[Payment]{I: *configs.GetCollection(Payments)}
