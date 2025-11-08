package utils

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims defines the custom JWT payload
type Claims struct {
	UserID uint64 `json:"uid"`   
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTManager handles generation and verification of JWTs
type JWTManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	accessExp  time.Duration
}

// NewJWTManager loads RSA keys and returns a JWT manager
func NewJWTManager(privPath, pubPath string, accessExp time.Duration) (*JWTManager, error) {
	privBytes, err := ioutil.ReadFile(privPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}

	pubBytes, err := ioutil.ReadFile(pubPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %v", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	return &JWTManager{
		privateKey: privKey,
		publicKey:  pubKey,
		accessExp:  accessExp,
	}, nil
}

// GenerateToken creates a signed JWT for a given user ID and role
func (jm *JWTManager) GenerateToken(userID uint64, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jm.accessExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(jm.privateKey)
}

// VerifyToken validates and decodes a JWT, returning its claims
func (jm *JWTManager) VerifyToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jm.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parse error: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
