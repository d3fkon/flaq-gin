package models

import (
	"github.com/d3fkon/gin-flaq/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuizTemplate struct {
	Id        primitive.ObjectID `bson:"_id"`
	Title     string             `bson:"Title" json:"Title" binding:"required"`
	Questions []Question         `bson:"Questions" json:"Questions" binding:"required"`
	CreatedAt primitive.DateTime `bson:"CreatedAt" json:"-" `
}

type Question struct {
	Question    string   `bson:"Question" json:"Question" binding:"required"`
	Description string   `bson:"Description" json:"Description" binding:"required"`
	Options     []string `bson:"Options" json:"Options" binding:"required"`
	AnswerIndex int      `bson:"AnswerIndex" json:"AnswerIndex" binding:"required"`
}

type QuizEntry struct {
	Id         primitive.ObjectID  `bson:"_id"`
	Campaign   primitive.ObjectID  `bson:"Campaign"`
	Template   primitive.ObjectID  `bson:"Template"`
	User       primitive.ObjectID  `bson:"User"`
	IsComplete bool                `bson:"IsComplete"`
	CreatedAt  primitive.Timestamp `bson:"CreatedAt"`
	UpdateAt   primitive.Timestamp `bson:"UpdatedAt"`
}

var QuziTemplateModel = Collection[QuizTemplate]{I: *configs.GetCollection(QuizTemplates)}
var QuizEntryModel = Collection[QuizEntry]{I: *configs.GetCollection(QuizEntries)}
