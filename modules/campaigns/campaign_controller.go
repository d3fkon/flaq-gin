package campaigns

import (
	"log"

	"github.com/d3fkon/gin-flaq/middleware"
	"github.com/d3fkon/gin-flaq/models"
	"github.com/d3fkon/gin-flaq/modules"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	modules.Controller
}

func Setup(g *gin.Engine) {
	c := Controller{}
	router := g.Group("/campaign")
	{
		router.POST("/a/", c.createCampaign)
		router.POST("/a/quiz/template", c.createQuizTemplate)
		router.POST("/a/quiz/add", c.addQuizToCampaign)
		authenticated := router.Group("/")
		authenticated.Use(middleware.UserAuth())
		{
			// Get all campaigns
			authenticated.GET("/", c.getAllCampaignsForUser)
			// Get a quiz template for a campaign
			authenticated.GET("/:campaignId/quiz/", c.getQuizTemplateForCampaign)
			// Get all participants
			authenticated.POST("/participate", c.participate)
			authenticated.POST("/quiz/evaluate", c.evaluateQuiz)
			authenticated.GET("/conversion", c.conversion)
		}
	}
}

type campaignBody struct {
	Title          string   `json:"Title"`
	Description    string   `json:"Description"`
	TickerName     string   `json:"TickerName" binding:"required"`
	TickerImgUrl   string   `json:"TickerImgUrl" binding:"required"`
	RequiredFlaq   int      `json:"RequiredFlaq" binding:"required"`
	FlaqReward     int      `json:"FlaqReward" binding:"required"`
	AirdropPerUser float64  `json:"AirdropPerUser" binding:"required"`
	TotalAirdrop   float64  `json:"TotalAirdrop" binding:"required"`
	CurrentAirdrop float64  `json:"CurrentAirdrop" binding:"required"`
	YTVideoUrl     string   `json:"YTVideoUrl" binding:"required"`
	ArticleUrls    []string `json:"ArticleUrls" binding:"required"`
}

// Create a campaign
// @Router   /campaign/a [post]
// @Summary  Create a campaign
// @Tags     Campaign Admin
// @Accept   application/json
// @Param    campaignBody  body  campaignBody  true  "Campaign Details"
func (c Controller) createCampaign(ctx *gin.Context) {
	body := campaignBody{}
	c.BindBody(ctx, &body)
	log.Println(body)
	res := CreateCampaign(&body)
	c.HandleResponse(ctx, res)
}

// Create a quiz template to be used
// @Router   /campaign/a/quiz/template [post]
// @Summary  Create a quiz template
// @Tags     Campaign Admin
// @Accept   application/json
// @Param    models.QuizTemplate  body  models.QuizTemplate  true  "Campaign Details"
func (c Controller) createQuizTemplate(ctx *gin.Context) {
	body := models.QuizTemplate{}
	c.BindBody(ctx, &body)
	CreateQuizTemplate(&body)
	c.HandleResponse(ctx, body)
}

//----------------------------

type addQuizToCampaignBody struct {
	CampaignId     string `json:"CampaignId"`
	QuizTemplateId string `json:"QuizTemplateId"`
}

// Create a quiz template to be used
// @Router   /campaign/a/quiz/add [post]
// @Summary  Create a quiz template [FOR ADMIN]
// @Tags     Campaign Admin
// @Accept   application/json
// @Param    addQuizToCampaignBody  body  addQuizToCampaignBody  true  "Quiz and Campaign Details"
func (c Controller) addQuizToCampaign(ctx *gin.Context) {
	body := addQuizToCampaignBody{}
	c.BindBody(ctx, &body)
	res := AddQuizToCampaign(body.CampaignId, body.QuizTemplateId)
	c.HandleResponse(ctx, res)
}

// Get all campaigns for a user
// @Router   /campaign/ [get]
// @Summary  Get all campaigns
// @param    Authorization  header  string  true  "Authorization"
// @Tags     Campaigns
// @Accept   application/json
func (c Controller) getAllCampaignsForUser(ctx *gin.Context) {
	user := c.ReqUser(ctx)
	res := GetCampaignParticipationForUser(user)
	c.HandleResponse(ctx, res)
}

// Get a quiz template for a campaign
// @Router   /campaign/{id}/quiz [get]
// @Summary  Get all campaigns
// @param    Authorization  header  string  true  "Authorization" "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImFzaHdpbkBvbnBhci5pbiIsImV4cCI6MTY1NzY1NjUzOX0.YhL-tGczj3OVlXNjCLwNJL4rkLdz7IQey4JFgmKmgoY"
// @param    id             path    string  true  "Campaign ID"
// @Tags     Campaigns
// @Accept   application/json
func (c Controller) getQuizTemplateForCampaign(ctx *gin.Context) {
	campaignId := ctx.Param("campaignId")
	res := GetQuizTemplateForCampaign(campaignId)
	c.HandleResponse(ctx, res)
}

// ----------------------------------------------
type quizParticipationBody struct {
	CampaignParticipationId string `json:"CampaignParticipationId" binding:"required"`
	QuizTemplateId          string `json:"QuizTemplateId" binding:"required"`
	Answers                 []int  `json:"Answers" binding:"required"`
}

// Evaluate a quiz
// @Router   /campaign/quiz/evaluate [post]
// @Summary  Evaluate a quiz submission
// @param    Authorization  header  string  true  "Authorization"
// @Tags     Campaigns
// @Accept   application/json
// @Param    quizParticipationBody  body  quizParticipationBody  true  "Quiz and Campaign Details"
func (c Controller) evaluateQuiz(ctx *gin.Context) {
	body := quizParticipationBody{}
	user := c.ReqUser(ctx)
	c.BindBody(ctx, &body)
	res := EvaluateQuiz(&user, body.CampaignParticipationId, body.QuizTemplateId, body.Answers)
	c.HandleResponse(ctx, res)
}

// -----------------------------------------------

type campaignParticipationBody struct {
	CampaignId string `json:"CampaignId"`
}

// Participate in campaign
// @Router   /campaign/participate [post]
// @Summary  Create a campaign
// @Tags     Campaigns
// @param    Authorization  header  string  true  "Authorization"
// @Accept   application/json
// @Param    campaignParticipationBody  body  campaignParticipationBody  true  "Campaign ID"
func (c Controller) participate(ctx *gin.Context) {
	body := campaignParticipationBody{}
	user := c.ReqUser(ctx)
	c.BindBody(ctx, &body)
	res := ParticipateInCampaign(body.CampaignId, user)
	c.HandleResponse(ctx, res)
}

type ticker struct {
	TickerName string  `json:"TickerName"`
	Conversion float64 `json:"Conversion"`
	Name       string  `json:"Name"`
}

// Convert a Ticker name to INR
// @Router   /campaign/conversion [get]
// @Summary  Get all conversions
// @param    Authorization  header  string  true  "Authorization" "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImFzaHdpbkBvbnBhci5pbiIsImV4cCI6MTY1NzY1NjUzOX0.YhL-tGczj3OVlXNjCLwNJL4rkLdz7IQey4JFgmKmgoY"
// @Tags     Conversion
// @Accept   application/json
func (c Controller) conversion(ctx *gin.Context) {
	tickers := []ticker{
		{
			TickerName: "USDT",
			Name:       "USD Tether",
			Conversion: 80.0,
		},
		{
			TickerName: "BNB",
			Name:       "Binance Token",
			Conversion: 19864.0,
		},
		{
			TickerName: "MATIC",
			Name:       "Polygon",
			Conversion: 59.23,
		},
	}
	c.ReqUser(ctx)
	c.HandleResponse(ctx, tickers)
}
