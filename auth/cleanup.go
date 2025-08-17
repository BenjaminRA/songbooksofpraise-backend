package auth

import (
	"time"

	"github.com/BenjaminRA/himnario-backend/db/sqlite"
)

// CleanupExpiredTokens removes expired tokens from the database
func CleanupExpiredTokens() error {
	db := sqlite.GetDBConnection()
	currentTime := time.Now().Unix()

	// Clean up expired session tokens
	_, err := db.Exec("DELETE FROM session_tokens WHERE at_exp < ? OR rt_exp < ?", currentTime, currentTime)
	if err != nil {
		return err
	}

	// Clean up expired verification tokens
	_, err = db.Exec("DELETE FROM verification_tokens WHERE expiration < ?", currentTime)
	if err != nil {
		return err
	}

	// Clean up expired reset tokens
	_, err = db.Exec("DELETE FROM reset_tokens WHERE expiration < ?", currentTime)
	if err != nil {
		return err
	}

	return nil
}

// RevokeUserTokens removes all tokens for a specific user
func RevokeUserTokens(userID int) error {
	db := sqlite.GetDBConnection()

	// Remove session tokens
	_, err := db.Exec("DELETE FROM session_tokens WHERE user_id = ?", userID)
	if err != nil {
		return err
	}

	// Remove verification tokens
	_, err = db.Exec("DELETE FROM verification_tokens WHERE user_id = ?", userID)
	if err != nil {
		return err
	}

	// Remove reset tokens
	_, err = db.Exec("DELETE FROM reset_tokens WHERE user_id = ?", userID)
	if err != nil {
		return err
	}

	return nil
}
