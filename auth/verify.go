package auth

import (
	"fmt"
	"os"

	"github.com/BenjaminRA/himnario-backend/db/sqlite"
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
			db := sqlite.GetDBConnection()

			var stored_user_id int
			err := db.QueryRow("SELECT user_id FROM session_tokens WHERE access_uuid = ?", sessionTokenClaims["access_uuid"].(string)).Scan(&stored_user_id)

			// Token belongs to different user or doesn't exist
			if err != nil || stored_user_id != int(sessionTokenClaims["user_id"].(float64)) {
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

		db := sqlite.GetDBConnection()
		var stored_user_id int
		err = db.QueryRow("SELECT user_id FROM session_tokens WHERE refresh_uuid = ?", refreshTokenClaims["refresh_uuid"].(string)).Scan(&stored_user_id)

		// Refresh Token Expired or doesn't exist
		if err != nil {
			return fmt.Errorf("refresh_token.expired")
		}

		// Token belongs to different user
		if stored_user_id != int(refreshTokenClaims["user_id"].(float64)) {
			return fmt.Errorf("session_token.missing")
		}

		user, err := new(models.User).GetUserById(stored_user_id)
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
