package campaigns

import (
	"fmt"
	"time"

	"github.com/d3fkon/gin-flaq/models"
	"github.com/d3fkon/gin-flaq/modules/users"
	"github.com/d3fkon/gin-flaq/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create a campaign
func CreateCampaign(campaign *models.Campaign) any {
	campaign.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	campaign.QuizIds = []primitive.ObjectID{}
	campaign.TaskType = models.TaskTypes.QUIZ
	campaign.Id = primitive.NewObjectID()
	models.CampaignModel.New(*campaign)
	return campaign
}

// Get all campaigns
func GetAllCampaigns(_ models.User) any {
	campaigns := []models.Campaign{}
	models.CampaignModel.FindMany(bson.M{}, &campaigns)
	return campaigns
}

// Craete quiz template
func CreateQuizTemplate(quizTemplate *models.QuizTemplate) any {
	quizTemplate.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	quizTemplate.Id = primitive.NewObjectID()
	fmt.Println(quizTemplate.Questions)
	models.QuziTemplateModel.New(*quizTemplate)
	return quizTemplate
}

// Get all campaign participations for user
func GetCampaignParticipationForUser(user models.User) any {
	query := bson.D{{
		Key: "UserId", Value: models.ObjId(user.Id.Hex()),
	}}

	// All campaigns the user is participating in
	participations := []models.CampaignParticipation{}
	if err := models.CampaignParticipationModel.FindMany(query, &participations); err != nil {
		utils.Panic(401, "Error finding user participations", err)
	}

	// All campaigns
	campaigns := []models.Campaign{}
	if err := models.CampaignModel.FindMany(bson.M{}, &campaigns); err != nil {
		utils.Panic(401, "Error finding campaigns", err)
	}

	return gin.H{
		"Participations": participations,
		"Campaigns":      campaigns,
	}
}

// Record participation in a campaign
func ParticipateInCampaign(campaignId string, user models.User) any {
	campaign := models.Campaign{}
	if err := models.CampaignModel.FindOneById(campaignId, &campaign); err != nil {
		utils.Panic(404, "Cannot find campaign", err)
	}
	requiredFlaq := campaign.RequiredFlaq
	if user.FlaqPoints < float64(requiredFlaq) {
		utils.Panic(401, "Low Flaq Point Balance", nil)
	}
	participation := models.CampaignParticipation{
		CreatedAt:  primitive.NewDateTimeFromTime(time.Now()),
		UserId:     models.ObjId(user.Id.Hex()),
		CampaignId: models.ObjId(campaign.Id.Hex()),
		IsComplete: false,
		FlaqSpent:  campaign.RequiredFlaq,
	}
	models.CampaignParticipationModel.New(participation)
	users.UpdateFlaqPoints(&user, -float64(requiredFlaq))
	return participation
}
