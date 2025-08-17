package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/BenjaminRA/himnario-backend/db/sqlite"
	"github.com/BenjaminRA/himnario-backend/helpers"
	"github.com/BenjaminRA/himnario-backend/models"
	"github.com/dgrijalva/jwt-go"
)

func ResetToken(user models.User) (string, error) {
	expiration := time.Now().Add(time.Hour * 24).Unix()
	key := helpers.HashValue(fmt.Sprintf("%dRESET%s", user.ID, os.Getenv("SECRET")))

	// Creating Reset token
	tokenClaims := jwt.MapClaims{}
	tokenClaims["user_id"] = user.ID

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	resetToken, err := at.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	// Saving token in SQLite
	db := sqlite.GetDBConnection()
	_, err = db.Exec("INSERT OR REPLACE INTO reset_tokens (key, token, user_id, expiration) VALUES (?, ?, ?, ?)",
		key, resetToken, user.ID, expiration)
	if err != nil {
		return "", err
	}

	return resetToken, nil
}

func VerifyResetToken(resetToken string) (models.User, error) {
	// Parsing the reset token
	token, err := jwt.Parse(resetToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("session_token.invalid")
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return models.User{}, err
	}

	// Validating the token
	tokenClaims := token.Claims.(jwt.MapClaims)
	user_id := int(tokenClaims["user_id"].(float64))
	user, err := new(models.User).GetUserById(user_id)
	if err != nil {
		return models.User{}, err
	}

	key := helpers.HashValue(fmt.Sprintf("%dRESET%s", user.ID, os.Getenv("SECRET")))

	// fetching token from SQLite
	db := sqlite.GetDBConnection()
	var storedToken string
	var expiration int64
	err = db.QueryRow("SELECT token, expiration FROM reset_tokens WHERE key = ?", key).Scan(&storedToken, &expiration)
	if err != nil {
		return models.User{}, err
	}

	// Check if token has expired
	if time.Now().Unix() > expiration {
		return models.User{}, fmt.Errorf("reset_token.expired")
	}

	// Check if tokens match
	if storedToken != resetToken {
		return models.User{}, fmt.Errorf("reset_token.invalid")
	}

	// Remove the reset token after successful verification
	_, err = db.Exec("DELETE FROM reset_tokens WHERE key = ?", key)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
