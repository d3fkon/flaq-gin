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
	M modules.Controller
}

func Setup(g *gin.Engine) {
	c := Controller{M: modules.Controller{}}
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

// Create User godoc
// @Router    /users/apply-referral [post]
// @Security  ApiKeyAuth
// @param     Authorization  header  string  true  "Authorization"
// @Summary   Apply a referral code
// @Tags      Users
// @Accept    application/json
// @Param     ApplyReferralBody  body  ApplyReferralBody  true  "Add Referral Data"
// @Produce   json
func (c Controller) ApplyReferral(ctx *gin.Context) {
	user := c.M.ReqUser(ctx)
	body := ApplyReferralBody{}
	c.M.BindBody(*ctx, &body)
	res := ApplyReferral(user, body.ReferralCode)
	c.M.HandleResponse(ctx, res)
}

// Get User godoc
// @Router    /users/profile [get]
// @Security  ApiKeyAuth
// @param     Authorization  header  string  true  "Authorization"
// @Summary   Get user profile
// @Tags      Users
// @Accept    application/json
// @Produce   json
func (c Controller) GetProfile(ctx *gin.Context) {
	user := c.M.ReqUser(ctx)
	c.M.HandleResponse(ctx, user)
}
