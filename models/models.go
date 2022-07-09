package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	Users    = "users"
	Payments = "payments"
)

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

type Collection struct {
	I mongo.Collection // Instance of the collection
}

func (c Collection) FindOneById(id string, model interface{}) error {
	ctx, cancel := GetContext()
	defer cancel()
	if err := c.I.FindOne(ctx, bson.M{"_id": id}).Decode(model); err != nil {
		return errors.New("Cannot find document")
	}
	return nil
}

func (c Collection) FindByIdAndUpdate(idHex string, update bson.M, updated interface{}) error {
	ctx, cancel := GetContext()
	defer cancel()
	if err := c.I.FindOneAndUpdate(ctx, bson.M{"_id": ObjId(idHex)}, update).Decode(updated); err != nil {
		fmt.Println(err)
		return errors.New("Cannot update document")
	}
	return nil
}

func (c Collection) FindOneAndUpdate(find bson.M, update bson.M, updated interface{}) error {
	ctx, cancel := GetContext()
	defer cancel()
	if err := c.I.FindOneAndUpdate(ctx, find, update).Decode(updated); err != nil {
		fmt.Println(err)
		return errors.New("Cannot update document")
	}
	return nil
}

func (c Collection) New(document interface{}) error {
	ctx, cancel := GetContext()
	defer cancel()
	_, err := c.I.InsertOne(ctx, document)
	return err
}

func (c Collection) FindMany(bson bson.M, elem interface{}) error {
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
