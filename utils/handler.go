package utils

import (
	"errors"
	"net/http"

	"github.com/d3fkon/gin-flaq/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"Message"`
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	}
	// Add other Messages here
	return "Unknown error"
}

func BindBody(ctx gin.Context, body interface{}) bool {
	if err := ctx.ShouldBindJSON(&body); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{fe.Field(), getErrorMsg(fe)}
			}
			Panic(http.StatusBadRequest, "Input fields are invalid", out)
		}
		return true
	}

	return false
}

// Panic struct
type HttpMessagePanic struct {
	StatusCode uint16 `json:"StatusCode"`
	Message    string `json:"Message"`
}

// Panic struct with errors
type HttpErrorPanic struct {
	StatusCode uint16      `json:"StatusCode"`
	Message    string      `json:"Message"`
	Errors     interface{} `json:"errors"`
}

// Helper function to create planned panics
func Panic(StatusCode uint16, Message string, errors interface{}) {
	if errors == nil {
		panic(HttpMessagePanic{StatusCode, Message})
	}
	panic(HttpErrorPanic{StatusCode, Message, errors})
}

// A custom recovery handler to handle panics
func HandleRecovery(ctx *gin.Context, recovered interface{}) {

	// Handle planned errors
	if err, ok := recovered.(HttpErrorPanic); ok {
		ctx.JSON(int(err.StatusCode), gin.H{
			"StatusCode": err.StatusCode,
			"Message":    err.Message,
			"Errors":     err.Errors,
		})
		return
	}

	// Handle planned errors
	if err, ok := recovered.(HttpMessagePanic); ok {
		ctx.JSON(int(err.StatusCode), gin.H{
			"StatusCode": err.StatusCode,
			"Message":    err.Message,
		})
		return
	}

	// Handle unplanned errors
	if err, ok := recovered.(string); ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": http.StatusInternalServerError,
			"Message":    err,
		})
		return
	}

	// Handle completely unplanned errors
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"StatusCode": http.StatusInternalServerError,
		"Message":    "Internal Server Error",
	})
}

// A wrapper function to handle all JSON responses
func HandleResponse(ctx *gin.Context, response interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": http.StatusOK,
		"Data":       response,
	})
}

// return just an object id
func ObjId(s string) primitive.ObjectID {
	o, _ := primitive.ObjectIDFromHex(s)
	return o
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ReqUser(ctx *gin.Context) models.User {
	user, exists := ctx.Get("user")
	if !exists {
		Panic(http.StatusUnauthorized, "Not authenticated", nil)
	}
	return user.(models.User)
}
