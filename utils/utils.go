package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

// GenerateVerificationToken generates a random token for email verification
func GenerateVerificationToken() string {
	bytes := make([]byte, 16) // Create a slice to hold 16 random bytes
	// Fill the slice with random bytes
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err) // If there's an error generating random bytes, log a fatal error
	}
	return hex.EncodeToString(bytes) // Convert the bytes to a hexadecimal string and return it
}
