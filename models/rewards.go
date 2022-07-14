package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// We are grouping all the rewards by Ticker Name
type Reward struct {
	Id                     primitive.ObjectID         `bson:"_id"`
	CampaignParticipations CampaignParticipationSlice `bson:"CampaignParticipations"`
	User                   UserWrapper                `bson:"User"`
	Amount                 float64                    `bson:"Amount"`
	TickerName             string                     `bson:"TickerName"`
	TickerImageUrl         string                     `bson:"TickerImageUrl" json:"TickerImageUrl"`
	CreatedAt              primitive.DateTime         `bson:"CreatedAt"`
}

type CampaignParticipationSlice struct {
	Ids  []primitive.ObjectID     `bson:"Ids"`
	Data *[]CampaignParticipation `bson:"Data"`
}

type RewardWrapper struct {
	Id   primitive.ObjectID `bson:"Id"`
	Data *Reward            `bson:"Data"`
}

var RewardModel = makeModel[Reward](Rewards)
