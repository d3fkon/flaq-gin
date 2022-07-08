package campaigns

import (
	"github.com/d3fkon/gin-flaq/middleware"
	"github.com/d3fkon/gin-flaq/modules"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	M modules.Controller
}

func Setup(g *gin.Engine) {
	c := Controller{M: modules.Controller{}}
	router := g.Group("/campaigns")
	router.Use(middleware.UserAuth())
	{
		router.POST("/create/template", c.createQuizTemplate)
		router.POST("/create/campaign", c.createQuizForCampaign)
		router.POST("/evaluate", c.evaluateQuiz)
		router.GET("/", c.getAllCampaignsForUser)
	}
}

// Create a quiz template to be used
func (c Controller) createQuizTemplate(ctx *gin.Context) {}

// Get all campaigns for a user
func (c Controller) getAllCampaignsForUser(ctx *gin.Context) {}

// Create a quiz for a user, for a campaign
func (c Controller) createQuizForCampaign(ctx *gin.Context) {}

// Evaluate a given quiz
func (c Controller) evaluateQuiz(ctx *gin.Context) {}
