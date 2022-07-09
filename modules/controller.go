package modules

import (
	"errors"
	"net/http"

	"github.com/d3fkon/gin-flaq/models"
	"github.com/d3fkon/gin-flaq/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller struct{}

func (c Controller) ReqUser(ctx *gin.Context) models.User {
	user, exists := ctx.Get("user")
	if !exists {
		utils.Panic(http.StatusUnauthorized, "Not authenticated", nil)
	}
	return user.(models.User)
}

// A wrapper function to handle all JSON responses
func (c Controller) HandleResponse(ctx *gin.Context, response interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": http.StatusOK,
		"Data":       response,
	})
}

func (c Controller) BindBody(ctx *gin.Context, body interface{}) bool {
	if err := ctx.ShouldBindJSON(&body); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]utils.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = utils.ErrorMsg{Field: fe.Field(), Message: utils.GetErrorMsg(fe)}
			}
			utils.Panic(http.StatusBadRequest, "Input fields are invalid", out)
		}
		return true
	}

	return false
}
