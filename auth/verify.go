package auth

import (
	"fmt"
	"os"

	redisdb "github.com/BenjaminRA/himnario-backend/db/redis"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func VerifyToken(c *gin.Context) error {
	missingSessionToken := false
	expiredSessionToken := false

	// Verify Session Token exists
	sessionTokenCookie, err := c.Request.Cookie("SessionToken")
	missingSessionToken = err != nil

	if !missingSessionToken {
		// Verify Session Token is valid
		sessionToken, err := jwt.Parse(sessionTokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("session_token.invalid")
			}

			return []byte(os.Getenv("ACCESS_SECRET")), nil
		})
		expiredSessionToken = err != nil

		if !expiredSessionToken {
			sessionTokenClaims := sessionToken.Claims.(jwt.MapClaims)
			redisdb := redisdb.GetRedisConnection()
			result := redisdb.Get(sessionTokenClaims["access_uuid"].(string))

			// Token belongs to different user
			if result.Err() == nil && result.Val() != sessionTokenClaims["user_id"] {
				return fmt.Errorf("session_token.invalid")
			}
		}

	}

	// Checking Refresh Token
	if missingSessionToken || expiredSessionToken {
		refreshTokenCookie, err := c.Request.Cookie("RefreshToken")
		if err != nil {
			return fmt.Errorf("refresh_token.missing")
		}

		refreshToken, err := jwt.Parse(refreshTokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("refresh_token.invalid")
			}

			return []byte(os.Getenv("REFRESH_SECRET")), nil
		})
		if err != nil {
			return fmt.Errorf("refresh_token.expired")
		}

		refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)

		redisdb := redisdb.GetRedisConnection()
		result := redisdb.Get(refreshTokenClaims["refresh_uuid"].(string))

		// Refresh Token Expired
		if result.Err() != nil {
			return fmt.Errorf("refresh_token.expired")
		}

		// Token belongs to different user
		if result.Val() != refreshTokenClaims["user_id"] {
			return fmt.Errorf("session_token.missing")
		}

		user, err := new(models.User).GetUserById(result.Val())
		if err != nil {
			return fmt.Errorf("refresh_token.invalid")
		}

		token, err := CreateToken(user)
		if err != nil {
			return fmt.Errorf("refresh_token.invalid")
		}

		token.SendToken(c)
		return nil
	}

	return nil
}
