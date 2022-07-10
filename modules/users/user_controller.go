package users

import (
	"github.com/d3fkon/gin-flaq/middleware"
	"github.com/d3fkon/gin-flaq/modules"
	"github.com/gin-gonic/gin"
)

type CreateUserBody struct {
	Email    string `binding:"required,email" json:"Email"`
	Password string `binding:"required" json:"Password"`
}

type Controller struct {
	modules.Controller
}

func Setup(g *gin.Engine) {
	c := Controller{}
	router := g.Group("/users")
	{
		authenticated := router.Group("/")
		{
			authenticated.Use(middleware.UserAuth())
			authenticated.POST("/apply-referral", c.ApplyReferral)
			authenticated.GET("/profile", c.GetProfile)
		}
	}
}

type ApplyReferralBody struct {
	ReferralCode string `binding:"required" json:"ReferralCode"`
}

// @Router    /users/apply-referral [post]
// @param     Authorization  header  string  true  "Authorization"
// @Summary   Apply a referral code
// @Tags      Users
// @Accept    application/json
// @Param     ApplyReferralBody  body  ApplyReferralBody  true  "Add Referral Data"
func (c Controller) ApplyReferral(ctx *gin.Context) {
	user := c.ReqUser(ctx)
	body := ApplyReferralBody{}
	c.BindBody(ctx, &body)
	res := ApplyReferral(user, body.ReferralCode)
	c.HandleResponse(ctx, res)
}

// @Router    /users/profile [get]
// @param     Authorization  header  string  true  "Authorization"
// @Summary   Get user profile
// @Tags      Users
// @Accept    application/json
// @Produce   json
func (c Controller) GetProfile(ctx *gin.Context) {
	user := c.ReqUser(ctx)
	c.HandleResponse(ctx, user)
}
