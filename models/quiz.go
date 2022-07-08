package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type QuizTemplate struct {
	Title     string             `bson:"Title"`
	Questions []Question         `bson:"Questions"`
	Campaign  primitive.ObjectID `bson:"Campaign"`
}

type Question struct {
	Question    string   `bson:"Question"`
	Description string   `bson:"Description"`
	Options     []string `bson:"Options"`
}

type QuizEntry struct {
	Campaign   primitive.ObjectID  `bson:"Campaign"`
	Template   QuizTemplate        `bson:"Template"`
	User       primitive.ObjectID  `bson:"User"`
	IsComplete bool                `bson:"IsComplete"`
	CreatedAt  primitive.Timestamp `bson:"CreatedAt"`
	UpdateAt   primitive.Timestamp `bson:"UpdatedAt"`
}
