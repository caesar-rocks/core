package core

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateAppKey generates a random 32 byte key for AES-256 encryption
func GenerateAppKey() string {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return ""
	}

	return hex.EncodeToString(key)
}
