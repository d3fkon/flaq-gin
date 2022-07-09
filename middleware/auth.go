package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/d3fkon/gin-flaq/jwt"
	"github.com/d3fkon/gin-flaq/models"
	"github.com/gin-gonic/gin"
)

const (
	AuthHeaderKey = "authorization"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get(AuthHeaderKey)

		user := models.User{}

		invalidTokenResponse := gin.H{
			"StatusCode": http.StatusUnauthorized,
			"Message":    "Invalid Access Token",
		}

		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, invalidTokenResponse)
			return
		}
		if len(strings.Split(accessToken, " ")) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, invalidTokenResponse)
			return
		}

		accessToken = strings.Split(accessToken, " ")[1]

		if err := (jwt.Jwt{}.ValidateAccessToken(accessToken, &user)); err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, invalidTokenResponse)
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
