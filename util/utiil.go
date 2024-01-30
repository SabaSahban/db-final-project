package util

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GenerateTrackingCode() string {
	// Implement a method to generate a unique tracking code for each transaction
	// You can use a combination of timestamp, random numbers, and account information
	// For simplicity, you can use a simple random code generation method
	rand.Seed(time.Now().UnixNano())
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	codeLength := 10
	code := make([]byte, codeLength)
	for i := 0; i < codeLength; i++ {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func HashPassword(password string) (string, error) {
	// Generate a hashed password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
