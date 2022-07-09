package models

import (
	"errors"
	"fmt"

	"github.com/d3fkon/gin-flaq/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id               primitive.ObjectID `bson:"_id"`
	Email            string             `bind:"email,required" bson:"Email"`
	FlaqPoints       float64            `bson:"FlaqPoints"`
	RewardMultiplier int                `bson:"RewardMultiplier"`
	ReferralCode     string             `bson:"ReferralCode"`
	IsAllowed        bool               `bson:"IsAllowed"`
	DeviceToken      string             `bson:"DeviceToken" json:"-"`
	RefreshToken     string             `bson:"RefreshToken" json:"-"`
	PasswordHash     string             `bson:"PasswordHash" json:"-"`
	ReferralData     Referral           `bson:"ReferralData" json:"-"`
	WalletAddresses  Wallet             `bson:"WalletAddresses"`
	CreatedAt        primitive.DateTime `bson:"CreatedAt"`
}

type Referral struct {
	AppliedReferral string               `bson:"AppliedReferral"`
	ReferredUsers   []primitive.ObjectID `bson:"ReferredUsers"`
}

type Wallet struct {
	Solana   string `bson:"Solana"`
	Ethereum string `bson:"Ethereum"`
	Avax     string `bson:"Avax"`
}

var referralIndex = configs.CreateIndex(Users, "ReferralCode", true, false)
var emailIndex = configs.CreateIndex(Users, "Email", true, false)

var UserModel = Collection{I: *configs.GetCollection(Users)}

func (c Collection) GetUserByEmail(email string, user *User) error {
	ctx, cancel := GetContext()
	defer cancel()
	if err := UserModel.I.FindOne(ctx, bson.M{"Email": email}).Decode(&user); err != nil {
		fmt.Println(err)
		return errors.New("Cannot find user")
	}
	return nil
}
