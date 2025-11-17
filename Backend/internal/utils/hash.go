package Utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	memory      = 64 * 1024
	iterations  = 3
	parallelism = 2
	saltLength  = 16
	keyLength   = 32
)

func HashPassword(password string) (string, error) {
	// generate random salt
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// Hash
	hash := argon2.IDKey([]byte(password), salt, iterations, memory, uint8(parallelism), keyLength)

	// Combine salt + hash in BASE64
	saltB64 := base64.StdEncoding.EncodeToString(salt)
	hashB64 := base64.StdEncoding.EncodeToString(hash)

	return saltB64 + ":" + hashB64, nil
}

// VerifyPassword compares raw password with stored salt:hash
func VerifyPassword(password, stored string) (bool, error) {
	parts := strings.Split(stored, ":")
	if len(parts) != 2 {
		return false, errors.New("invalid stored password format")
	}

	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}

	storedHash, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	// Hash input password using extracted salt
	hash := argon2.IDKey([]byte(password), salt, iterations, memory, uint8(parallelism), keyLength)

	// Compare
	if len(hash) != len(storedHash) {
		return false, nil
	}

	// Constant-time compare
	isMatch := true
	for i := 0; i < len(hash); i++ {
		if hash[i] != storedHash[i] {
			isMatch = false
		}
	}

	return isMatch, nil
}
