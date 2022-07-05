package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"Message"`
}

func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "numeric":
		return "Should be a numeric"
	}
	// Add other Messages here
	return "Unknown error"
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
