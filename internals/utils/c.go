package utils

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

type CryptoUtils struct{}

func (c CryptoUtils) Hash(plaintext string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (c CryptoUtils) Match(plaintext, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext))
}

func (c CryptoUtils) GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

var Crypto = CryptoUtils{}
