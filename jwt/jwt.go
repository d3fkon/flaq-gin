package jwt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/d3fkon/gin-flaq/configs"
	"github.com/d3fkon/gin-flaq/models"
	"github.com/d3fkon/gin-flaq/utils"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
)

type Jwt struct {
}

var userModel = configs.GetCollection(models.Users)

// A get user by email redundant helper to prevent circular deps
func getUserByEmail(email string, user *models.User) error {
	ctx, cancel := utils.GetContext()
	defer cancel()

	if err := userModel.FindOne(ctx, bson.M{"Email": email}).Decode(&user); err != nil {
		return errors.New("Cannot find user")
	}
	return nil
}

func (j Jwt) CreateToken(user models.User) (models.Token, error) {
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
		if err := getUserByEmail(email, user); err != nil {
			return errors.New("Cannot find user")
		}
		return nil
	}

	return errors.New("invalid token")
}

func (Jwt) ValidateRefreshToken(model models.Token, user *models.User) error {
	sha1 := sha1.New()
	io.WriteString(sha1, configs.GetEnv(configs.JWT_SECRET))

	salt := string(sha1.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	data, err := base64.URLEncoding.DecodeString(model.RefreshToken)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	if string(plain) != model.AccessToken {
		return errors.New("invalid token")
	}

	claims := jwt.MapClaims{}
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(model.AccessToken, claims)

	if err != nil {
		return err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token")
	}

	email := payload["Email"].(string)

	if err = getUserByEmail(email, user); err != nil {
		return errors.New("Cannot find user")
	}

	return nil
}

func (Jwt) createRefreshToken(token models.Token) (models.Token, error) {
	sha1 := sha1.New()
	io.WriteString(sha1, configs.GetEnv(configs.JWT_SECRET))

	salt := string(sha1.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		fmt.Println(err.Error())

		return token, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return token, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return token, err
	}

	token.RefreshToken = base64.URLEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(token.AccessToken), nil))

	return token, nil
}
