package campaigns

import (
	"fmt"
	"time"

	"github.com/d3fkon/gin-flaq/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create a campaign
func CreateCampaign(campaign *models.Campaign) any {
	campaign.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	campaign.Id = primitive.NewObjectID()
	models.CampaignModel.New(*campaign)
	return campaign
}

// Craete quiz template
func CreateQuizTemplate(quizTemplate *models.QuizTemplate) any {
	quizTemplate.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	quizTemplate.Id = primitive.NewObjectID()
	fmt.Println(quizTemplate.Questions)
	models.QuziTemplateModel.New(*quizTemplate)
	return quizTemplate
}
