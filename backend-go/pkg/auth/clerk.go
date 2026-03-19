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

// ValidateClerkToken validates a Clerk JWT token and returns the claims.
// Uses the Decode → GetJSONWebKey → Verify pattern per Clerk SDK v2 docs.
// clerk.SetKey() must be called before this function is used.
func ValidateClerkToken(tokenString string) (*ClerkClaims, error) {
	ctx := context.Background()

	// Step 1: Decode the token to extract the Key ID
	decoded, err := jwt.Decode(ctx, &jwt.DecodeParams{
		Token: tokenString,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to decode Clerk token: %w", err)
	}

	// Step 2: Fetch the JSON Web Key for this token's Key ID
	jwk, err := jwt.GetJSONWebKey(ctx, &jwt.GetJSONWebKeyParams{
		KeyID: decoded.KeyID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get Clerk JWK: %w", err)
	}

	// Step 3: Verify the token with the JWK
	claims, err := jwt.Verify(ctx, &jwt.VerifyParams{
		Token: tokenString,
		JWK:   jwk,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to verify Clerk token: %w", err)
	}

	return &ClerkClaims{SessionClaims: *claims}, nil
}

// GetUsername extracts the Clerk user ID from claims.
// Note: claims.Subject is the Clerk user ID (e.g. "user_2abc..."), not a username.
func GetUsername(claims *ClerkClaims) string {
	return claims.Subject
}
