package campaigns

import (
	"fmt"
	"log"
	"time"

	"github.com/d3fkon/gin-flaq/models"
	"github.com/d3fkon/gin-flaq/modules/rewards"
	"github.com/d3fkon/gin-flaq/modules/users"
	"github.com/d3fkon/gin-flaq/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create a campaign for a from an admin
// Creates a new campaign in the database
func CreateCampaign(data *campaignBody) models.Campaign {
	campaign := models.Campaign{
		CreatedAt: models.Now(),
		Quizzes: models.QuizSliceWrapper{
			Ids: []primitive.ObjectID{},
		},
		Title:          data.Title,
		Description:    data.Description,
		TickerName:     data.TickerName,
		TickerImageUrl: data.TickerImgUrl,
		TaskType:       models.TaskTypes.QUIZ,
		FlaqReward:     data.FlaqReward,
		CurrentAirdrop: data.CurrentAirdrop,
		TotalAirdrop:   data.TotalAirdrop,
		AirdropPerUser: data.AirdropPerUser,
		RequiredFlaq:   data.RequiredFlaq,
		Id:             primitive.NewObjectID(),
		YTVideoUrl:     data.YTVideoUrl,
		ArticleUrls:    data.ArticleUrls,
	}

	if err := models.CampaignModel.New(campaign); err != nil {
		utils.Panic(400, "Error occured creating", err)
	}
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
	models.QuizTemplateModel.New(*quizTemplate)
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
			"Quizzes.Ids": models.ObjId(quizTemplateId),
		},
	}
	updated := models.Campaign{}
	if err := models.CampaignModel.FindOneAndUpdate(query, update, &updated); err != nil {
		utils.Panic(400, "Campagin not found", err)
	}
	return updated
}

// GetCampaignParticipationForUser will return a Map, which has the user's participations
// The map contains a) All campaigns, b) Campaigns participated by the given user
// The service panics if there is an error finding campaign participations or while populating it
func GetCampaignParticipationForUser(user models.User) gin.H {
	query := bson.D{{
		Key: "User.Id", Value: models.ObjId(user.Id.Hex()),
	}}

	// All campaigns the user is participating in
	participations := []models.CampaignParticipation{}
	if err := models.CampaignParticipationModel.FindMany(query, &participations); err != nil {
		utils.Panic(400, "Error finding user participations", err)
	}

	// All campaigns
	campaigns := []models.Campaign{}
	if err := models.CampaignModel.FindMany(bson.M{}, &campaigns); err != nil {
		utils.Panic(400, "Error finding campaigns", err)
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
		utils.Panic(400, "Low Flaq Point Balance", nil)
	}
	campaignParticipation := models.CampaignParticipation{
		Id:        primitive.NewObjectID(),
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		User: models.UserWrapper{
			Id: models.ObjId(user.Id.Hex()),
		},
		Campaign: models.CampaignWrapper{
			Id: models.ObjId(campaign.Id.Hex()),
		},
		IsComplete: false,
		FlaqSpent:  campaign.RequiredFlaq,
	}
	models.CampaignParticipationModel.New(campaignParticipation)
	users.UpdateFlaqPoints(&user, -float64(requiredFlaq))
	return campaignParticipation
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
		As:           "Quizzes.Data",
		ForeignModel: models.QuizTemplates,
		LocalField:   "Quizzes.Ids",
	}

	campaigns := []models.Campaign{}
	if err := models.CampaignModel.FindManyPopulate(query, populate, &campaigns); err != nil {
		utils.Panic(400, "[1] Campaign Not Found", err)
	}
	if len(campaigns) > 0 && len(*&campaigns[0].Quizzes.Ids) > 0 {
		quizzes := campaigns[0].Quizzes
		return &(*quizzes.Data)[0]
	}
	utils.Panic(400, "[2] Campaign Not Found", nil)
	return nil
}

// Check if the Campaign ID corelates to a Campaign
// Check if the QuizTemplateId corelates to a Quiz Template
// Panic if any of the above two fail
// Evaluate the answers by checking if the answer indexes stored and the answers match
// Create a QuizEntry with the above created data
func EvaluateQuiz(user *models.User, campaignParticipationId, quizTemplateId string, answers []int) models.QuizEntry {
	campaignParticipation := models.CampaignParticipation{}
	campaign := models.Campaign{}
	quizTemplate := models.QuizTemplate{}
	err1 := models.CampaignParticipationModel.FindOneById(campaignParticipationId, &campaignParticipation)
	err2 := models.QuizTemplateModel.FindOneById(quizTemplateId, &quizTemplate)
	err3 := models.CampaignModel.FindOneById(campaignParticipation.Campaign.Id.Hex(), &campaign)
	if err1 != nil || err2 != nil || err3 != nil {
		log.Printf("E1 %e\nE2 %e\nE3 %e", err1, err2, err3)
		utils.Panic(400, "Error evaluating participation", err1)
	}
	if len(answers) != len(quizTemplate.Questions) {
		error := fmt.Sprintf("Invalid answers array length. Expected %d got %d", len(quizTemplate.Questions), len(answers))
		utils.Panic(400, error, nil)
	}
	// Evaluate the quiz by checking if the indexes of answers and the data match
	score := 0
	for i, question := range quizTemplate.Questions {
		answer := answers[i]
		if answer == question.AnswerIndex {
			score++
		}
	}
	isQuizPassing := score == len(quizTemplate.Questions)
	// Check the number of question in the quiz
	// Check the number of correct answers in the quiz
	quizEntry := models.QuizEntry{
		IsPassing:     isQuizPassing,
		CreatedAt:     models.Now(),
		Id:            primitive.NewObjectID(),
		QuestionCount: len(quizTemplate.Questions),
		CorrectCount:  score,
	}

	quizEntry.Quiz.Id = models.ObjId(quizTemplateId)
	quizEntry.Campaign.Id = models.ObjId(campaignParticipation.Campaign.Id.Hex())
	quizEntry.User.Id = models.ObjId(user.Id.Hex())
	if err := models.QuizEntryModel.New(quizEntry); err != nil {
		utils.Panic(400, "Error creating quiz", err)
	}

	// If the task type was quiz and the quiz is answered successfully
	// Trigger a campaign completion event
	// Currently only handles task types which are QUIZ. Any other task type should be handled here
	if isQuizPassing {
		// USER IS PASSING THE QUIZ
		// Ideally all this should be an event sent to Kafka
		update := bson.M{
			"$set": bson.M{
				"IsComplete": true,
			},
		}
		updated := models.CampaignParticipation{}
		models.CampaignParticipationModel.FindByIdAndUpdate(campaignParticipationId, update, &updated)
		users.UpdateFlaqPoints(user, float64(campaign.FlaqReward))
		rewards.AddRewardsByParticipation(user, &campaignParticipation, &campaign)
		// TODO: Give crypto reward
	}

	return quizEntry
}
