package auth

import (
	"context"
	"fmt"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
)

// ClerkClaims represents the claims from a Clerk JWT
type ClerkClaims struct {
	clerk.SessionClaims
}

// ValidateClerkToken validates a Clerk JWT token and returns the claims
func ValidateClerkToken(tokenString string) (*ClerkClaims, error) {
	// Verify the token using Clerk's JWT verification
	// Clerk will automatically fetch the JWKS from their servers for verification
	// This works in both keyless mode and production
	ctx := context.Background()
	claims, err := jwt.Verify(ctx, &jwt.VerifyParams{
		Token: tokenString,
		// No JWK or JWKSClient needed - Clerk SDK will use default backend
		// to fetch public keys for verification
	})
	if err != nil {
		return nil, fmt.Errorf("failed to verify Clerk token: %w", err)
	}

	return &ClerkClaims{SessionClaims: *claims}, nil
}

// GetUsername extracts username from Clerk claims
// Clerk stores username in the session claims
func GetUsername(claims *ClerkClaims) string {
	// Clerk provides the user ID in the Subject field
	return claims.Subject
}
