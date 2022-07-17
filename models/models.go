package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/d3fkon/gin-flaq/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	Users                  = "users"
	Payments               = "payments"
	QuizTemplates          = "quiz_templates"
	QuizEntries            = "quiz_entries"
	Campaigns              = "campaigns"
	CampaignParticipations = "campaign_participations"
	Rewards                = "rewards"
)

type Models interface{}

// Cannot use utils due to circular dependency
func GetContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx, cancel
}

// return just an object id
func ObjId(s string) primitive.ObjectID {
	o, _ := primitive.ObjectIDFromHex(s)
	return o
}

func Now() primitive.DateTime {
	return primitive.NewDateTimeFromTime(time.Now())
}

type Collection[A Models] struct {
	I    mongo.Collection // Instance of the collection
	name string
}

func makeModel[A Models](name string) *Collection[A] {
	return &Collection[A]{
		I:    *configs.GetCollection(name),
		name: name,
	}
}

func (c Collection[M]) FindOne(query bson.M, model *M) error {
	ctx, cancel := GetContext()
	defer cancel()
	if err := c.I.FindOne(ctx, query).Decode(model); err != nil {
		return errors.New("cannot find document")
	}
	return nil
}

func (c Collection[M]) FindOneById(id string, model *M) error {
	ctx, cancel := GetContext()
	defer cancel()
	if err := c.I.FindOne(ctx, bson.M{"_id": ObjId(id)}).Decode(model); err != nil {
		return errors.New("Cannot find document")
	}
	return nil
}

func (c Collection[M]) FindByIdAndUpdate(idHex string, update bson.M, updated *M) error {
	ctx, cancel := GetContext()
	defer cancel()
	if err := c.I.FindOneAndUpdate(ctx, bson.M{"_id": ObjId(idHex)}, update).Decode(updated); err != nil {
		fmt.Println(err)
		return errors.New("Cannot update document")
	}
	return nil
}

func (c Collection[M]) FindOneAndUpdate(find bson.M, update bson.M, updated *M) error {
	ctx, cancel := GetContext()
	defer cancel()
	if err := c.I.FindOneAndUpdate(ctx, find, update).Decode(updated); err != nil {
		fmt.Println(err)
		return errors.New("Cannot update document")
	}
	return nil
}

func (c Collection[M]) New(document M) error {
	ctx, cancel := GetContext()
	defer cancel()
	_, err := c.I.InsertOne(ctx, document)
	return err
}

// Find many. Basically send a bson.M or bson.D query
func (c Collection[M]) FindMany(bson interface{}, elem *[]M) error {
	ctx, cancel := GetContext()
	defer cancel()
	cursor, err := c.I.Find(ctx, bson)
	if err != nil {
		return err
	}
	if err := cursor.All(ctx, elem); err != nil {
		return err
	}
	return nil
}

type Populate struct {
	LocalField, ForeignModel, As string
}

func (c Collection[M]) FindManyPopulate(matchQuery bson.D, populate Populate, elem *[]M) error {
	ctx, cancel := GetContext()
	defer cancel()

	match := bson.D{{Key: "$match", Value: matchQuery}}
	lookup := bson.D{{Key: "$lookup", Value: bson.D{{
		Key:   "from",
		Value: populate.ForeignModel,
	}, {
		Key:   "localField",
		Value: populate.LocalField,
	}, {
		Key:   "foreignField",
		Value: "_id",
	}, {
		Key:   "as",
		Value: populate.As,
	}}}}

	// unwind := bson.D{{
	// 	Key: "$unwind",
	// 	Value: bson.D{{
	// 		Key:   "path",
	// 		Value: strings.Join([]string{"$", populate.As}, ""),
	// 	}, {
	// 		Key:   "preserveNullAndEmptyArrays",
	// 		Value: false,
	// 	}},
	// }}

	cursor, err := c.I.Aggregate(ctx, mongo.Pipeline{match, lookup})
	if err != nil {
		return err
	}
	if err := cursor.All(ctx, elem); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
