package rewards

import (
	"github.com/d3fkon/gin-flaq/middleware"
	"github.com/d3fkon/gin-flaq/modules"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	modules.Controller
}

func Setup(g *gin.Engine) {
	c := Controller{}
	router := g.Group("/rewards")
	router.Use(middleware.UserAuth())
	{
		router.GET("/", c.getAllRewards) // Get all rewards by grouping them by ticker id
		router.GET("/claim")             // Claim a particular reward
	}
}

// @Router   /rewards [get]
// @Summary  Get all rewards for a user
// @param    Authorization  header  string  true  "Authorization"
// @Tags     Rewards
// @Accept   application/json
// @Produce  json
// Get all rewards for a user
func (c Controller) getAllRewards(ctx *gin.Context) {
	user := c.ReqUser(ctx)
	res := GetAllRewards(&user)
	c.HandleResponse(ctx, res)
}
