package models

import (
	"time"

	"github.com/d3fkon/gin-flaq/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	Id        primitive.ObjectID `bson:"_id"`
	User      primitive.ObjectID `bson:"User" binding:"required"`
	Amount    string             `bson:"Amount" binding:"required"`
	CreatedAt time.Time          `bson:"CreatedAt"`
}

var PaymentModel = Collection{I: *configs.GetCollection(Payments)}
