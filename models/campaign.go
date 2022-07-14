package models

import (
	"github.com/d3fkon/gin-flaq/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskTypes struct {
	QUIZ string
}

var TaskTypes taskTypes = taskTypes{
	QUIZ: "QUIZ",
}

type Campaign struct {
	Id             primitive.ObjectID   `bson:"_id" json:"Id"`
	QuizIds        []primitive.ObjectID `bson:"QuizIds" json:"QuizIds"`
	Quizzes        *[]QuizTemplate      `bson:"Quizzes" json:"Quizzes"`
	Name           string               `bson:"Name" json:"Name"`
	Description    string               `bson:"Description" json:"Description" binding:"required"`
	RequiredFlaq   int                  `bson:"RequiredFlaq" json:"RequiredFlaq" binding:"required"`
	FlaqReward     int                  `bson:"FlaqReward" json:"FlaqReward"`
	TickerName     string               `bson:"TickerName" json:"TickerName" binding:"required"`
	TickerImageUrl string               `bson:"TickerImageUrl" json:"TickerImageUrl" binding:"required"`
	AirdropPerUser float64              `bson:"AirdropPerUser" json:"AirdropPerUser" binding:"required"`
	TotalAirdrop   float64              `bson:"TotalAirdrop" json:"TotalAirdrop" binding:"required"`
	CurrentAirdrop float64              `bson:"CurrentAirdrop" json:"CurrentAirdrop" binding:"required"`
	CreatedAt      primitive.DateTime   `bson:"CreatedAt" json:"-"`
	TaskType       string               `bson:"TaskType" json:"TaskType"`
	ArticleUrls    []string             `bson:"ArticleUrls" json:"ArticleUrls"`
	YTVideoUrl     string               `bson:"YTVideoUrl" json:"YTVideoUrl"`
}

type CampaignWrapper struct {
	Id   primitive.ObjectID `bson:"Id"`
	Data *Campaign          `bson:"Data"`
}

var CampaignModel = Collection[Campaign]{I: *configs.GetCollection(Campaigns)}

type CampaignParticipation struct {
	Id         primitive.ObjectID `bson:"_id" json:"Id"`
	Campaign   CampaignWrapper    `bson:"Campaign" json:"Campaign"`
	User       UserWrapper        `bson:"User" json:"User"`
	CreatedAt  primitive.DateTime `bson:"CreatedAt" json:"CreatedAt"`
	IsComplete bool               `bson:"IsComplete" json:"IsComplete"`
	FlaqSpent  int                `bson:"FlaqSpent" json:"FlaqSpent"`
}

var CampaignParticipationModel = makeModel[CampaignParticipation](CampaignParticipations)
