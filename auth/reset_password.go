package auth

import (
	"fmt"
	"os"
	"time"

	redisdb "github.com/BenjaminRA/himnario-backend/db/redis"
	"github.com/BenjaminRA/himnario-backend/helpers"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/dgrijalva/jwt-go"
)

func ResetToken(user models.User) (string, error) {
	expiration := time.Now().Add(time.Hour * 24).Unix()
	key := helpers.HashValue(fmt.Sprintf("%sRESET%s", user.ID.Hex(), os.Getenv("SECRET")))

	// Creating Verification token
	tokenClaims := jwt.MapClaims{}
	tokenClaims["user_id"] = user.ID.Hex()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	resetToken, err := at.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	// Saving token in redis
	client := redisdb.GetRedisConnection()
	client.Set(key, resetToken, time.Until(time.Unix(expiration, 0)))

	return resetToken, nil
}

func VerifyResetToken(resetToken string) (models.User, error) {
	token, err := jwt.Parse(resetToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("session_token.invalid")
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return models.User{}, err
	}

	tokenClaims := token.Claims.(jwt.MapClaims)
	user, err := new(models.User).GetUserById(tokenClaims["tokenClaims"].(string))
	if err != nil {
		return models.User{}, err
	}

	key := helpers.HashValue(fmt.Sprintf("%sVERIFICATION%s", user.ID.Hex(), os.Getenv("SECRET")))

	// fetching token from redis
	client := redisdb.GetRedisConnection()
	result := client.Get(key)

	if result.Err() != nil {
		return models.User{}, result.Err()
	}

	return user, nil
}
