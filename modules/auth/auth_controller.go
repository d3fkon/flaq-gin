package auth

import (
	"net/http"

	"github.com/d3fkon/gin-flaq/jwt"
	"github.com/d3fkon/gin-flaq/models"
	"github.com/d3fkon/gin-flaq/modules/users"
	"github.com/d3fkon/gin-flaq/utils"
	"github.com/gin-gonic/gin"
)

type Controller struct{}

func (c Controller) Setup(g *gin.Engine) {
	router := g.Group("/auth")
	{
		router.POST("/signup", c.Signup)
		router.POST("/login", c.Login)
	}
}

type SignupBody struct {
	Email    string `binding:"required,email" json:"Email"`
	Password string `binding:"required,min=6,max=16" json:"Password"`
}

// Create User godoc
// @Router   /auth/signup [post]
// @Summary  User signup
// @Tags     Auth
// @Accept   application/json
// @Param    SignupBody  body  SignupBody  true  "Add Data"
// @Produce  json
func (c Controller) Signup(ctx *gin.Context) {
	body := SignupBody{}
	utils.BindBody(*ctx, &body)
	// TODO: Validate Password
	user, _ := users.CreateUser(users.CreateUserBody{
		Email:    body.Email,
		Password: body.Password,
	})
	token := genTokenAndSetCookie(ctx, &user)
	utils.HandleResponse(ctx, token)
}

type LoginBody struct {
	Email    string `binding:"required,email" json:"Email"`
	Password string `binding:"required" json:"Password"`
}

// User Login godoc
// @Router   /auth/login [post]
// @Summary  User login
// @Tags     Auth
// @Accept   application/json
// @Param    LoginBody  body  LoginBody  true  "Enter Login details"
// @Produce  json
func (c Controller) Login(ctx *gin.Context) {
	body := LoginBody{}
	utils.BindBody(*ctx, &body)
	user, isLoggedIn := users.CheckLogin(body.Email, body.Password)
	if !isLoggedIn {
		utils.Panic(http.StatusBadRequest, "Invalid Password", nil)
		return
	}
	token := genTokenAndSetCookie(ctx, &user)
	utils.HandleResponse(ctx, token)
}

func genTokenAndSetCookie(ctx *gin.Context, user *models.User) models.Token {
	token, err := jwt.Jwt{}.CreateToken(*user)
	if err != nil {
		utils.Panic(http.StatusInternalServerError, "Unable to create tokens", err)
	}
	users.UpdateRefreshToken(user, token.RefreshToken)
	// TODO: Set secure http only cookies in the future
	ctx.SetCookie("access-token", token.AccessToken, 1000000, "/", "0.0.0.0", false, false)
	ctx.SetCookie("refresh-token", token.RefreshToken, 1000000, "/", "0.0.0.0", false, false)
	return token
}
