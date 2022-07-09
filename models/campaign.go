package models

import (
	"github.com/d3fkon/gin-flaq/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Campaign struct {
	Quizzes        []primitive.ObjectID `bson:"Quizzes"`
	Name           string               `bson:"Name"`
	Description    string               `bson:"Description"`
	RequiredFlaq   int                  `bson:"RequiredFlaq"`
	TickerName     string               `bson:"TickerName"`
	TotalAirdrop   int                  `bson:"TotalAirdrop"`
	CurrentAirdrop int                  `bson:"CurrentAirdrop"`
}

var CampaignModel = Collection{I: *configs.GetCollection(Campaigns)}
