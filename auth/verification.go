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

func VerificationToken(user models.User) (string, error) {
	expiration := time.Now().Add(time.Hour * 24).Unix()
	key := helpers.HashValue(fmt.Sprintf("%sVERIFICATION%s", user.ID.Hex(), os.Getenv("SECRET")))

	// Creating Verification token
	tokenClaims := jwt.MapClaims{}
	tokenClaims["user_id"] = user.ID.Hex()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	verificationToken, err := at.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	// Saving token in redis
	client := redisdb.GetRedisConnection()
	client.Set(key, verificationToken, time.Until(time.Unix(expiration, 0)))

	return verificationToken, nil
}

func VerifyVerificationToken(verificationToken string) (models.User, error) {
	token, err := jwt.Parse(verificationToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("email.verify.invalid")
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return models.User{}, err
	}

	tokenClaims := token.Claims.(jwt.MapClaims)
	user, err := new(models.User).GetUserById(tokenClaims["user_id"].(string))
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

	user.Verified = true
	if err = user.UpdateUser(); err != nil {
		return models.User{}, err
	}

	return user, nil
}
