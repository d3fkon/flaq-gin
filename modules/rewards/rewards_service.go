package rewards

import (
	"log"

	"github.com/d3fkon/gin-flaq/models"
	"github.com/d3fkon/gin-flaq/modules/conversions"
	"github.com/d3fkon/gin-flaq/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Once a user completes a campaign, reward the user with the right reward by tracking it
// This function checks if the user already has a reward with the given ticker, if so, then updates it
// If there is not reward with the said ticker, then create a new reward document and initialize it
func AddRewardsByParticipation(user *models.User, participation *models.CampaignParticipation, campaign *models.Campaign) {
	reward := models.Reward{}
	query := bson.M{
		"TickerName": campaign.TickerName,
	}
	// Create a new reward object for the user if it doesn't exist
	if err := models.RewardModel.FindOne(query, &reward); err != nil {
		reward.TickerImageUrl = campaign.TickerImageUrl
		reward.TickerName = campaign.TickerName
		reward.CreatedAt = models.Now()
		reward.User.Id = user.Id
		reward.CampaignParticipations.Ids = []primitive.ObjectID{models.ObjId(participation.Id.Hex())}
		reward.Amount = campaign.AirdropPerUser
		reward.Id = primitive.NewObjectID()
		models.RewardModel.New(reward)
	} else {
		update := bson.M{
			"$set": bson.M{
				"Amount": reward.Amount + campaign.AirdropPerUser,
			},
			"$push": bson.M{
				"CampaignParticipations.Ids": models.ObjId(participation.Id.Hex()),
			},
		}
		models.RewardModel.FindByIdAndUpdate(reward.Id.Hex(), update, &reward)
	}
}

type rewardWithTicker struct {
	models.Reward
	conversions.Ticker
}

// Get all the rewards for a given user
// Get all the rewards which are already grouped by ticker name, populate and return
func GetAllRewards(user *models.User) []rewardWithTicker {
	rewards := []models.Reward{}
	query := bson.D{{
		Key:   "User.Id",
		Value: models.ObjId(user.Id.Hex()),
	}}
	populate := models.Populate{
		As:           "CampaignParticipations.Data",
		ForeignModel: models.CampaignParticipations,
		LocalField:   "CampaignParticipations.Ids",
	}

	if err := models.RewardModel.FindManyPopulate(query, populate, &rewards); err != nil {
		utils.Panic(400, "No rewards found for user", nil)
	}

	modifiedRewards := []rewardWithTicker{}
	for _, reward := range rewards {
		ticker, err := conversions.GetTickerByName(reward.TickerName)
		if err != nil {
			log.Println(err)
			continue
		}
		r := rewardWithTicker{
			reward,
			ticker,
		}
		modifiedRewards = append(modifiedRewards, r)
	}
	return modifiedRewards
}

// Claim rewards from the user's wallet
// This requires a smart contract call to transfer funds
func ClaimRewards(user *models.User) {
	// TODO: Make a smart contract call to send the right ticker
}
