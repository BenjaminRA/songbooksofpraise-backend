package models

import (
	"database/sql"
	"time"

	"github.com/BenjaminRA/himnario-backend/db/sqlite"
)

type PasswordResetToken struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

// CreatePasswordResetToken creates a new password reset token
func (prt *PasswordResetToken) Create() error {
	db := sqlite.GetDBConnection()

	result, err := db.Exec(`
		INSERT INTO password_reset_tokens (user_id, token, expires_at)
		VALUES (?, ?, ?)
	`, prt.UserID, prt.Token, prt.ExpiresAt)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	prt.ID = int(id)
	prt.CreatedAt = time.Now()
	prt.Used = false

	return nil
}

// GetByToken retrieves a password reset token by its token string
func (prt *PasswordResetToken) GetByToken(token string) error {
	db := sqlite.GetDBConnection()

	err := db.QueryRow(`
		SELECT id, user_id, token, expires_at, used, created_at
		FROM password_reset_tokens
		WHERE token = ? AND used = 0 AND expires_at > datetime('now')
	`, token).Scan(&prt.ID, &prt.UserID, &prt.Token, &prt.ExpiresAt, &prt.Used, &prt.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}

	return nil
}

// MarkAsUsed marks the token as used
func (prt *PasswordResetToken) MarkAsUsed() error {
	db := sqlite.GetDBConnection()

	_, err := db.Exec(`
		UPDATE password_reset_tokens
		SET used = 1
		WHERE id = ?
	`, prt.ID)

	return err
}

// InvalidateAllUserTokens marks all tokens for a user as used
func InvalidateAllUserTokens(userID int) error {
	db := sqlite.GetDBConnection()

	_, err := db.Exec(`
		UPDATE password_reset_tokens
		SET used = 1
		WHERE user_id = ? AND used = 0
	`, userID)

	return err
}

// DeleteExpiredTokens removes all expired tokens from the database
func DeleteExpiredTokens() error {
	db := sqlite.GetDBConnection()

	_, err := db.Exec(`
		DELETE FROM password_reset_tokens
		WHERE expires_at < datetime('now')
	`)

	return err
}
