package campaigns

import (
	"github.com/d3fkon/gin-flaq/middleware"
	"github.com/d3fkon/gin-flaq/models"
	"github.com/d3fkon/gin-flaq/modules"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	M modules.Controller
}

func Setup(g *gin.Engine) {
	c := Controller{M: modules.Controller{}}
	router := g.Group("/campaign")
	{
		router.POST("/", c.createCampaign)
		router.POST("/quiz/template", c.createQuizTemplate)
		router.POST("/quiz/create", c.createQuizForCampaign)
		router.POST("/evaluate", c.evaluateQuiz)
		authenticated := router.Group("/")
		authenticated.Use(middleware.UserAuth())
		authenticated.GET("/", c.getAllCampaignsForUser)
	}
}

// Create a campaign
// @Router    /campaign [post]
// @Summary   Create a campaign
// @Tags      Campaigns
// @Accept    application/json
// @Param     models.Campaign body  models.Campaign true  "Campaign Details"
func (c Controller) createCampaign(ctx *gin.Context) {
	body := models.Campaign{}
	c.M.BindBody(ctx, &body)
	CreateCampaign(&body)
	c.M.HandleResponse(ctx, body)
}

// Create a quiz template to be used
// @Router    /campaign/quiz/template [post]
// @Summary   Create a quiz template
// @Tags      Campaigns
// @Accept    application/json
// @Param     models.QuizTemplate body  models.QuizTemplate true  "Campaign Details"
func (c Controller) createQuizTemplate(ctx *gin.Context) {
	body := models.QuizTemplate{}
	c.M.BindBody(ctx, &body)
	CreateQuizTemplate(&body)
	c.M.HandleResponse(ctx, body)
}

// Get all campaigns for a user
// @Router    /campaign/ [get]
// @Summary   Get all campaigns
// @param    Authorization  header  string  true  "Authorization"
// @Tags      Campaigns
// @Accept    application/json
func (c Controller) getAllCampaignsForUser(ctx *gin.Context) {
	user := c.M.ReqUser(ctx)
	c.M.HandleResponse(ctx, GetAllCampaigns(user))
}

// Create a quiz for a user, for a campaign
func (c Controller) createQuizForCampaign(ctx *gin.Context) {}

// Evaluate a given quiz
func (c Controller) evaluateQuiz(ctx *gin.Context) {}
