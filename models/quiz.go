package models

import (
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

type QuizWrapper struct {
	Id   primitive.ObjectID `bson:"Id" json:"Id"`
	Data *QuizTemplate      `bson:"Data" json:"Data"`
}

type QuizSliceWrapper struct {
	Ids  []primitive.ObjectID `bson:"Ids" json:"Ids"`
	Data *[]QuizTemplate      `bson:"Data" json:"Data"`
}

type QuizEntry struct {
	Id            primitive.ObjectID `bson:"_id"`
	Campaign      CampaignWrapper    `bson:"Campaign"`
	Quiz          QuizWrapper        `bson:"Quiz"`
	User          UserWrapper        `bson:"User"`
	IsPassing     bool               `bson:"IsPassing"`
	QuestionCount int                `bson:"QuestionCount"`
	CorrectCount  int                `bson:"CorrectCount"`
	CreatedAt     primitive.DateTime `bson:"CreatedAt"`
}

var QuizTemplateModel = makeModel[QuizTemplate](QuizTemplates)
var QuizEntryModel = makeModel[QuizEntry](QuizEntries)
