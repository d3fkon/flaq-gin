package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/d3fkon/gin-flaq/configs"
	"github.com/d3fkon/gin-flaq/models"
	"github.com/golang-jwt/jwt"
)

type Jwt struct {
}

// Create access and refresh tokens for a user
func (j Jwt) CreateToken(user *models.User) (models.Token, error) {
	var err error

	claims := jwt.MapClaims{}
	claims["Email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt := models.Token{}

	jwt.AccessToken, err = token.SignedString([]byte(configs.GetEnv(configs.JWT_SECRET)))
	if err != nil {
		return jwt, err
	}

	return j.createRefreshToken(jwt)
}

func (Jwt) ValidateAccessToken(accessToken string, user *models.User) error {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(configs.GetEnv(configs.JWT_SECRET)), nil
	})

	if err != nil {
		return err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		email := payload["Email"].(string)
		if err := models.UserModel.GetUserByEmail(email, user); err != nil {
			return errors.New("Cannot find user")
		}
		return nil
	}

	return errors.New("invalid token")
}

func (Jwt) ValidateRefreshToken(refreshToken string, user *models.User) error {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return errors.New("invalid token")
	}

	claims := jwt.MapClaims{}
	parser := jwt.Parser{}
	token, _, err = parser.ParseUnverified(payload["token"].(string), claims)
	if err != nil {
		return err
	}

	payload, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token")
	}

	user.Email = payload["Email"].(string)

	return nil
}

func (Jwt) createRefreshToken(token models.Token) (models.Token, error) {
	var err error

	claims := jwt.MapClaims{}
	claims["token"] = token.AccessToken
	claims["exp"] = time.Now().Add(time.Hour * 24 * 90).Unix()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return token, err
	}

	return token, nil
}
