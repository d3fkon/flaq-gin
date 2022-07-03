package users

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/d3fkon/gin-flaq/configs"
	"github.com/d3fkon/gin-flaq/models"
	"github.com/d3fkon/gin-flaq/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var UserModel = configs.GetCollection(models.Users)

// Get user by Email
func GetUserByEmail(email string, user *models.User) error {
	ctx, cancel := utils.GetContext()
	defer cancel()

	if err := UserModel.FindOne(ctx, bson.M{"Email": email}).Decode(&user); err != nil {
		return errors.New("Cannot find user")
	}
	return nil
}

// Get user by Referral Code
func GetUserByReferralCode(referralCode string) (models.User, error) {
	ctx, cancel := utils.GetContext()
	defer cancel()

	referralUser := models.User{}
	if err := UserModel.FindOne(ctx, bson.M{"ReferralCode": referralCode}).Decode(&referralUser); err != nil {
		return referralUser, errors.New("Cannot find user")
	}
	return referralUser, nil
}

func CreateUser(data CreateUserBody) (models.User, error) {
	ctx, cancel := utils.GetContext()
	defer cancel()

	passwordHash, e := utils.HashPassword(data.Password)
	if e != nil {
		utils.Panic(http.StatusInternalServerError, "Error creating user [2]", e)
	}

	user := models.User{
		Email:            data.Email,
		Id:               primitive.NewObjectID(),
		IsAllowed:        false,
		RewardMultiplier: "1",
		FlaqPoints:       "0",
		PasswordHash:     passwordHash,
		ReferralCode:     utils.GenerateReferral(),
		ReferralData: models.Referral{
			ReferredUsers:   make([]primitive.ObjectID, 0),
			AppliedReferral: "",
		},
		WalletAddresses: models.Wallet{},
	}

	_, err := UserModel.InsertOne(ctx, user)
	if err != nil {
		utils.Panic(http.StatusInternalServerError, "Cannot create a new user", err)
	}

	return user, nil
}

func UpdateRefreshToken(user *models.User, refreshToken string) {
	ctx, cancel := utils.GetContext()
	defer cancel()
	user.RefreshToken = refreshToken
	if err := UserModel.FindOneAndUpdate(ctx, bson.M{"Email": user.Email}, bson.M{"$set": bson.M{"RefreshToken": refreshToken}}).Decode(&user); err != nil {
		utils.Panic(http.StatusInternalServerError, "Error Updating the Database", err)
	}
}

func CheckLogin(email string, password string) (models.User, bool) {
	user := models.User{}
	if err := GetUserByEmail(email, &user); err != nil {
		utils.Panic(http.StatusNotFound, "User not found", err)
	}
	return user, utils.CheckPasswordHash(password, user.PasswordHash)
}

func ApplyReferral(user models.User, referral string) interface{} {
	ctx, cancel := utils.GetContext()
	defer cancel()

	if user.IsAllowed {
		utils.Panic(http.StatusBadRequest, "Referrral code already applied", nil)
	}

	_, err := GetUserByReferralCode(referral)
	if err != nil {
		utils.Panic(http.StatusBadRequest, "Invalid referral code", nil)
	}

	rUpdateData := bson.M{
		"$push": bson.M{
			"ReferralData.ReferredUsers": user.Id,
		},
	}
	uUpdateData := bson.M{
		"$set": bson.M{
			"ReferralData.AppliedReferral": referral,
		},
	}
	var x interface{}
	if err := UserModel.FindOneAndUpdate(ctx, bson.M{"ReferralCode": referral}, rUpdateData).Decode(&x); err != nil {
		utils.Panic(http.StatusInternalServerError, "Error Updating the Database", err.Error())
	}
	if err := UserModel.FindOneAndUpdate(ctx, bson.M{"_id": utils.ObjId(user.Id.Hex())}, uUpdateData).Decode(&user); err != nil {
		fmt.Println(err)
		utils.Panic(http.StatusInternalServerError, "Error Updating the Database", err.Error())
	}

	return "Referral Applied"

}
