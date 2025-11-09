package Utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// HashPassword securely hashes a plain-text password using Argon2id.
// Argon2id is resistant to brute-force and GPU attacks, making it suitable for production.
func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)
	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 2, 32)
	return fmt.Sprintf("%s:%s", base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash)), nil
}

// function for verifying or matching password wiht salt in Aargon2
func VerifyPassword(password, encoded string) bool {
	var saltB64, hashB64 string
	fmt.Sscanf(encoded, "%[^:]:%s", &saltB64, &hashB64)
	salt, _ := base64.RawStdEncoding.DecodeString(saltB64)
	hash, _ := base64.RawStdEncoding.DecodeString(hashB64)
	test := argon2.IDKey([]byte(password), salt, 3, 64*1024, 2, 32)
	if len(test) != len(hash) {
		return false
	}
	var diff byte
	for i := range test {
		diff |= test[i] ^ hash[i]
	}
	return diff == 0
}
