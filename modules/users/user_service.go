package users

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/d3fkon/gin-flaq/models"
	"github.com/d3fkon/gin-flaq/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Get user by Referral Code
func GetUserByReferralCode(referralCode string) (models.User, error) {
	ctx, cancel := utils.GetContext()
	defer cancel()

	referralUser := models.User{}
	if err := models.UserModel.I.FindOne(ctx, bson.M{"ReferralCode": referralCode}).Decode(&referralUser); err != nil {
		return referralUser, errors.New("Cannot find user")
	}
	return referralUser, nil
}

func CreateUser(data CreateUserBody) (models.User, error) {
	passwordHash, e := utils.HashPassword(data.Password)
	if e != nil {
		utils.Panic(http.StatusInternalServerError, "Error creating user [2]", e)
	}

	user := models.User{
		Email:            data.Email,
		Id:               primitive.NewObjectID(),
		IsAllowed:        false,
		RewardMultiplier: 1,
		FlaqPoints:       0,
		PasswordHash:     passwordHash,
		ReferralCode:     utils.GenerateReferral(),
		ReferralData: models.Referral{
			ReferredUsers:   make([]primitive.ObjectID, 0),
			AppliedReferral: "",
		},
		WalletAddresses: models.Wallet{},
		CreatedAt:       primitive.NewDateTimeFromTime(time.Now()),
	}
	if err := models.UserModel.New(user); err != nil {
		utils.Panic(http.StatusInternalServerError, "Cannot create a new user", err)
	}

	return user, nil
}

// A helper method to update the user's flaq points balance
func UpdateFlaqPoints(user *models.User, reward int) {
	currentPoints := user.FlaqPoints
	update := bson.M{
		"$set": bson.M{
			"FlaqPoints": strconv.Itoa(currentPoints + reward),
		},
	}
	models.UserModel.FindByIdAndUpdate(user.Id.Hex(), update, &user)
}

func UpdateRefreshToken(user *models.User, refreshToken string) {
	user.RefreshToken = refreshToken
	if err := models.UserModel.FindOneAndUpdate(bson.M{"Email": user.Email}, bson.M{"$set": bson.M{"RefreshToken": refreshToken}}, &user); err != nil {
		utils.Panic(http.StatusInternalServerError, "Error Updating the Database", err)
	}
}

func CheckLogin(email string, password string) (models.User, bool) {
	user := models.User{}
	if err := models.UserModel.GetUserByEmail(email, &user); err != nil {
		utils.Panic(http.StatusNotFound, "User not found", err)
	}
	return user, utils.CheckPasswordHash(password, user.PasswordHash)
}

func ApplyReferral(user models.User, referral string) interface{} {
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
			"IsAllowed":                    true,
			"ReferralData.AppliedReferral": referral,
		},
	}
	var x interface{} // Empty dummy interface
	if err := models.UserModel.FindOneAndUpdate(bson.M{"ReferralCode": referral}, rUpdateData, &x); err != nil {
		utils.Panic(http.StatusInternalServerError, "[1] Error Updating the Database", err.Error())
	}
	if err := models.UserModel.FindByIdAndUpdate(user.Id.Hex(), uUpdateData, &user); err != nil {
		fmt.Println(err)
		utils.Panic(http.StatusInternalServerError, "[2] Error Updating the Database", err.Error())
	}

	return "Referral Applied"

}
