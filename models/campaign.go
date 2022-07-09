package models

import (
	"github.com/d3fkon/gin-flaq/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Campaign struct {
	Id             primitive.ObjectID   `bson:"_id" json:"_id"`
	Quizzes        []primitive.ObjectID `bson:"Quizzes" json:"Quizzes"`
	Name           string               `bson:"Name" json:"Name" binding:"required"`
	Description    string               `bson:"Description" json:"Description" binding:"required"`
	RequiredFlaq   int                  `bson:"RequiredFlaq" json:"RequiredFlaq" binding:"required"`
	TickerName     string               `bson:"TickerName" json:"TickerName" binding:"required"`
	TotalAirdrop   int                  `bson:"TotalAirdrop" json:"TotalAirdrop" binding:"required"`
	CurrentAirdrop int                  `bson:"CurrentAirdrop" json:"CurrentAirdrop" binding:"required"`
	CreatedAt      primitive.DateTime   `bson:"CreatedAt" json:"-"`
}

var CampaignModel = Collection[Campaign]{I: *configs.GetCollection(Campaigns)}
