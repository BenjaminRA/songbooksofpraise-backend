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

func VerificationToken(user models.User) (string, error) {
	expiration := time.Now().Add(time.Hour * 24).Unix()
	key := helpers.HashValue(fmt.Sprintf("%dVERIFICATION%s", user.ID, os.Getenv("SECRET")))

	// Creating Verification token
	tokenClaims := jwt.MapClaims{}
	tokenClaims["user_id"] = user.ID

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	verificationToken, err := at.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	// Saving token in SQLite
	db := sqlite.GetDBConnection()
	_, err = db.Exec("INSERT OR REPLACE INTO verification_tokens (key, token, user_id, expiration) VALUES (?, ?, ?, ?)",
		key, verificationToken, user.ID, expiration)
	if err != nil {
		return "", err
	}

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
	user_id := int(tokenClaims["user_id"].(float64))
	user, err := new(models.User).GetUserById(user_id)
	if err != nil {
		return models.User{}, err
	}

	key := helpers.HashValue(fmt.Sprintf("%dVERIFICATION%s", user.ID, os.Getenv("SECRET")))

	// fetching token from SQLite
	db := sqlite.GetDBConnection()
	var storedToken string
	var expiration int64
	err = db.QueryRow("SELECT token, expiration FROM verification_tokens WHERE key = ?", key).Scan(&storedToken, &expiration)
	if err != nil {
		return models.User{}, err
	}

	// Check if token has expired
	if time.Now().Unix() > expiration {
		return models.User{}, fmt.Errorf("email.verify.expired")
	}

	// Check if tokens match
	if storedToken != verificationToken {
		return models.User{}, fmt.Errorf("email.verify.invalid")
	}

	user.Verified = true
	if err = user.UpdateUser(); err != nil {
		return models.User{}, err
	}

	// Remove the verification token after successful verification
	_, err = db.Exec("DELETE FROM verification_tokens WHERE key = ?", key)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
