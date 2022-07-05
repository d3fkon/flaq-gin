package main

import (
	"fmt"
	"net/http"

	_ "github.com/d3fkon/gin-flaq/configs"
	_ "github.com/d3fkon/gin-flaq/docs"
	"github.com/gin-contrib/cors"

	"github.com/d3fkon/gin-flaq/modules/auth"
	"github.com/d3fkon/gin-flaq/modules/payments"
	"github.com/d3fkon/gin-flaq/modules/users"
	"github.com/d3fkon/gin-flaq/utils"
	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerFiles "github.com/swaggo/files"
)

// swagger embed files

// gin-swagger middleware
// swagger embed files

var db = make(map[string]string)

func setupRouter(r *gin.Engine) {
	// Disable Console Color
	// gin.DisableConsoleColor()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	auth.Setup(r)
	users.Setup(r)
	payments.Setup(r)
}

// @title           Flaq API
// @version         2.0
// @description     This is a sample server server.
// @termsOfService  http://flaq.club/privacy-policy
// @contact.name    Ashwin Prasad
// @contact.url     http://www.swagger.io/support
// @contact.email   ashwin@flaq.club
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
// @host            0.0.0.0:8080
// @BasePath        /
// @schemes         http
func main() {
	r := gin.Default()
	r.Use(gin.CustomRecovery(utils.HandleRecovery))
	r.Use(cors.Default())

	setupRouter(r)

	// Listen and Server in 0.0.0.0:8080
	fmt.Println("Running on http://0.0.0.0:8080")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
