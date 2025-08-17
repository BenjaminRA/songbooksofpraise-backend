package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/BenjaminRA/himnario-backend/helpers"
)

type AppToken struct {
	AppID     string `json:"app_id"`
	AppKey    string `json:"app_key"`
	ExpiresAt int64  `json:"expires_at"`
}

// createKey creates a 32-byte key from the secret string
func createKey(secret string) []byte {
	hash := sha256.Sum256([]byte(secret))
	return hash[:]
}

// decryptAES decrypts data using AES-256-GCM
func decryptAES(encryptedData []byte) ([]byte, error) {
	helpers.LoadLocalEnv()

	// Create key from environment variable
	secret := os.Getenv("APP_SECRET")
	key := createKey(secret)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Extract nonce and ciphertext
	nonceSize := aesGCM.NonceSize()

	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("invalid encrypted data")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]

	// Decrypt
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// ValidateAppToken verifies and decrypts the token
func ValidateAppToken(encryptedToken string) error {
	// Decode from base64
	encryptedData, err := base64.URLEncoding.DecodeString(encryptedToken)
	if err != nil {
		return fmt.Errorf("invalid token format: %w", err)
	}

	// Decrypt the token
	decryptedData, err := decryptAES(encryptedData)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	// Parse JSON
	var token AppToken
	if err := json.Unmarshal(decryptedData, &token); err != nil {
		return fmt.Errorf("invalid token structure")
	}

	// Validate app ID
	if token.AppID != "songbooksofpraise" {
		return fmt.Errorf("invalid app ID")
	}

	// Check if token has expired
	now := time.Now().Unix()
	if now > token.ExpiresAt {
		return fmt.Errorf("token has expired")
	}

	return nil
}
