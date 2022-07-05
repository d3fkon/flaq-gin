package payments

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
	router := g.Group("/payments")
	router.Use(middleware.UserAuth())
	{
		router.POST("/register", c.register)
		router.GET("/", c.getAll)
	}
}

type registerPaymentBody struct {
	Amount string `json:"Amount" binding:"required,numeric"`
}

// Register a payment for a user
// User Register a Payment godoc
// @Router   /payments/register [post]
// @Summary  Register a new UPI Payment
// @param    Authorization  header  string  true  "Authorization"
// @Tags     Payments
// @Accept   application/json
// @Param    registerPaymentBody  body  registerPaymentBody  true  "Register Payment"
// @Produce  json
func (c Controller) register(ctx *gin.Context) {
	user := c.M.ReqUser(ctx)
	body := registerPaymentBody{}
	c.M.BindBody(*ctx, &body)
	res := RegisterPayment(user, body.Amount)
	c.M.HandleResponse(ctx, res)
}

// Register a payment for a user
// User Get all Payments godoc
// @Router   /payments [get]
// @Summary  Register a new UPI Payment
// @param    Authorization  header  string  true  "Authorization"
// @Tags     Payments
// @Accept   application/json
// @Produce  json
// Get all payments for a user
func (c Controller) getAll(ctx *gin.Context) {
	user := c.M.ReqUser(ctx)
	c.M.HandleResponse(ctx, GetAllPayments(user))
}
