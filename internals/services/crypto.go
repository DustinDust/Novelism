package services

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

type CryptoService struct{}

func NewCryptoService() CryptoService {
	return CryptoService{}
}

func (service CryptoService) Hash(plaintext string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (service CryptoService) Match(plaintext, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext))
}

func (service CryptoService) GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
