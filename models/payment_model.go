package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	Id         primitive.ObjectID `bson:"_id"`
	User       primitive.ObjectID `bson:"User" binding:"required"`
	Amount     string             `bson:"Amount" binding:"required"`
	CreatedAt  primitive.DateTime `bson:"CreatedAt"`
	FlaqReward float64            `bson:"FlaqReward"`
}

var PaymentModel = makeModel[Payment](Payments)
