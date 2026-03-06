package auth

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT claims
type Claims struct {
	Username string `json:"sub"`
	Role     string `json:"rol"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(username, role string) (string, error) {
	// Get JWT secret from environment
	secretStr := os.Getenv("JWT_SECRET")
	if secretStr == "" {
		return "", fmt.Errorf("JWT_SECRET not set")
	}

	// Decode base64 secret
	secret, err := base64.StdEncoding.DecodeString(secretStr)
	if err != nil {
		return "", fmt.Errorf("invalid JWT_SECRET: %w", err)
	}

	// Get expiration hours (default 6)
	expirationHours := 6
	if expStr := os.Getenv("JWT_EXPIRATION_HOURS"); expStr != "" {
		if hours, err := strconv.Atoi(expStr); err == nil {
			expirationHours = hours
		}
	}

	// Create claims
	now := time.Now()
	claims := Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * time.Duration(expirationHours))),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*Claims, error) {
	// Get JWT secret from environment
	secretStr := os.Getenv("JWT_SECRET")
	if secretStr == "" {
		return nil, fmt.Errorf("JWT_SECRET not set")
	}

	// Decode base64 secret
	secret, err := base64.StdEncoding.DecodeString(secretStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_SECRET: %w", err)
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
