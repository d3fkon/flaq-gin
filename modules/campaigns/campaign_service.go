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

// Create a campaign for a from an admin
// Creates a new campaign in the database
func CreateCampaign(campaign *models.Campaign) any {
	campaign.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	campaign.QuizIds = []primitive.ObjectID{}
	campaign.TaskType = models.TaskTypes.QUIZ
	campaign.Id = primitive.NewObjectID()
	models.CampaignModel.New(*campaign)
	return campaign
}

// Get all campaigns
// Get all campaigns including the user's participating users
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

// AddQuizToCampaign is a helper service method which adds a given QuizTemplate with an ID into a Campaign
// This performs a very simple PUSH operation and doesn't expect anything more than the campaign ID and the
// Quiz Template ID
func AddQuizToCampaign(campaignId string, quizTemplateId string) models.Campaign {
	query := bson.M{
		"_id": models.ObjId(campaignId),
	}

	update := bson.M{
		"$push": bson.M{
			"QuizIds": models.ObjId(quizTemplateId),
		},
	}
	updated := models.Campaign{}
	if err := models.CampaignModel.FindOneAndUpdate(query, update, &updated); err != nil {
		utils.Panic(401, "Campagin not found", err)
	}
	return updated
}

// GetCampaignParticipationForUser will return a Map, which has the user's participations
// The map contains a) All campaigns, b) Campaigns participated by the given user
// The service panics if there is an error finding campaign participations or while populating it
func GetCampaignParticipationForUser(user models.User) gin.H {
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
// Check for the required flaq points to participate in a campaign
// Deduct the right amount of Flaq from the user for the same and enrol the user in the campaign
// TODO: Check for the campaign capacity
func ParticipateInCampaign(campaignId string, user models.User) models.CampaignParticipation {
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

// GetQuizTemplateForCampaign returns the most recent quiz template for a given campaign
// The campaign should be a legitimate campaign with quizzes
// The service will panic if the campaign does not exist or if the campaign does not have quizzes
func GetQuizTemplateForCampaign(campaignId string) *models.QuizTemplate {
	query := bson.D{{
		Key:   "_id",
		Value: models.ObjId(campaignId),
	}}

	populate := models.Populate{
		As:           "Quizzes",
		ForeignModel: models.QuizTemplates,
		LocalField:   "QuizIds",
	}

	campaigns := []models.Campaign{}
	if err := models.CampaignModel.FindManyPopulate(query, populate, &campaigns); err != nil {
		utils.Panic(401, "[1] Campaign Not Found", err)
	}
	if len(campaigns) > 0 && len(*campaigns[0].Quizzes) > 0 {
		quizzes := campaigns[0].Quizzes
		return &(*quizzes)[0]
	}
	utils.Panic(401, "[2] Campaign Not Found", nil)
	return nil
}
