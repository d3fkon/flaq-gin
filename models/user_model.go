package models

import (
	"github.com/d3fkon/gin-flaq/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id               primitive.ObjectID `bson:"_id"`
	Email            string             `bind:"email,required" bson:"Email"`
	FlaqPoints       string             `bson:"FlaqPoints"`
	RewardMultiplier string             `bson:"RewardMultiplier"`
	ReferralCode     string             `bson:"ReferralCode"`
	IsAllowed        bool               `bson:"IsAllowed"`
	DeviceToken      string             `bson:"DeviceToken" json:"-"`
	RefreshToken     string             `bson:"RefreshToken"`
	PasswordHash     string             `bson:"PasswordHash" json:"-"`
	ReferralData     Referral           `bson:"ReferralData" json:"-"`
	WalletAddresses  Wallet             `bson:"WalletAddresses"`
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
