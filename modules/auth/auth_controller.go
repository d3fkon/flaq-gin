package auth

import (
	"net/http"

	"github.com/d3fkon/gin-flaq/jwt"
	"github.com/d3fkon/gin-flaq/models"
	"github.com/d3fkon/gin-flaq/modules"
	"github.com/d3fkon/gin-flaq/modules/users"
	"github.com/d3fkon/gin-flaq/utils"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	modules.Controller
}

func Setup(g *gin.Engine) {
	c := Controller{}
	router := g.Group("/auth")
	{
		router.POST("/signup", c.signup)
		router.POST("/login", c.login)
		router.POST("/token/refresh", c.getAccessToken)
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
func (c Controller) signup(ctx *gin.Context) {
	body := SignupBody{}
	c.BindBody(ctx, &body)
	// TODO: Validate Password
	user, _ := users.CreateUser(users.CreateUserBody{
		Email:    body.Email,
		Password: body.Password,
	})
	token := genTokenAndSetCookie(ctx, &user)
	c.HandleResponse(ctx, token)
}

type LoginBody struct {
	Email    string `binding:"required,email" json:"Email"`
	Password string `binding:"required" json:"Password"`
}

// User login godoc
// @Router   /auth/login [post]
// @Summary  User login
// @Tags     Auth
// @Accept   application/json
// @Param    LoginBody  body  LoginBody  true  "Enter login details"
// @Produce  json
func (c Controller) login(ctx *gin.Context) {
	body := LoginBody{}
	c.BindBody(ctx, &body)
	user, isLoggedIn := users.CheckLogin(body.Email, body.Password)
	if !isLoggedIn {
		utils.Panic(http.StatusBadRequest, "Invalid Password", nil)
		return
	}
	token := genTokenAndSetCookie(ctx, &user)
	c.HandleResponse(ctx, token)
}

type RefreshTokenBody struct {
	RefreshToken string `binding:"required" json:"RefreshToken"`
}

// User refresh-token godoc
// @Router   /auth/token/refresh [post]
// @Summary  Issue a new Access token
// @Tags     Auth
// @Accept   application/json
// @Param    RefreshTokenBody  body  RefreshTokenBody  true  "Refresh token"
// @Produce  json
func (c Controller) getAccessToken(ctx *gin.Context) {
	// Craete the body
	body := RefreshTokenBody{}
	user := models.User{}
	c.BindBody(ctx, &body)
	jwt := jwt.Jwt{}

	// Validate the refresh token
	if err := (jwt.ValidateRefreshToken(body.RefreshToken, &user)); err != nil {
		utils.Panic(http.StatusUnauthorized, "Refresh token invalid", nil)
	}

	token := genTokenAndSetCookie(ctx, &user)

	// If the token is valid, then generate a new set of tokens
	c.HandleResponse(ctx, token)
}

func genTokenAndSetCookie(ctx *gin.Context, user *models.User) models.Token {
	token, err := jwt.Jwt{}.CreateToken(user)
	if err != nil {
		utils.Panic(http.StatusInternalServerError, "Unable to create tokens", err)
	}
	users.UpdateRefreshToken(user, token.RefreshToken)
	// TODO: Set secure http only cookies in the future
	ctx.SetCookie("access-token", token.AccessToken, 1000000, "/", "0.0.0.0", false, false)
	ctx.SetCookie("refresh-token", token.RefreshToken, 1000000, "/", "0.0.0.0", false, false)
	return token
}
