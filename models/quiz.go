package models

import (
	"github.com/d3fkon/gin-flaq/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuizTemplate struct {
	Title     string     `bson:"Title"`
	Questions []Question `bson:"Questions"`
}

type Question struct {
	Question    string   `bson:"Question"`
	Description string   `bson:"Description"`
	Options     []string `bson:"Options"`
}

type QuizEntry struct {
	Campaign   primitive.ObjectID  `bson:"Campaign"`
	Template   primitive.ObjectID  `bson:"Template"`
	User       primitive.ObjectID  `bson:"User"`
	IsComplete bool                `bson:"IsComplete"`
	CreatedAt  primitive.Timestamp `bson:"CreatedAt"`
	UpdateAt   primitive.Timestamp `bson:"UpdatedAt"`
}

var QuziTemplateModel = Collection{I: *configs.GetCollection(QuizTemplates)}
var QuizEntryModel = Collection{I: *configs.GetCollection(QuizEntries)}
