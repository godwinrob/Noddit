package auth

import (
	"encoding/base64"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testSecret is a valid base64-encoded secret for testing
var testSecret = base64.StdEncoding.EncodeToString([]byte("this-is-a-test-secret-key-32bytes!"))

func setupJWTEnv(t *testing.T) {
	t.Helper()
	t.Setenv("JWT_SECRET", testSecret)
}

func TestGenerateToken_Success(t *testing.T) {
	setupJWTEnv(t)

	token, err := GenerateToken("testuser", "user")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateToken_MissingSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	_, err := GenerateToken("testuser", "user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_SECRET not set")
}

func TestGenerateToken_InvalidBase64Secret(t *testing.T) {
	t.Setenv("JWT_SECRET", "%%%not-valid-base64%%%")

	_, err := GenerateToken("testuser", "user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid JWT_SECRET")
}

func TestValidateToken_RoundTrip(t *testing.T) {
	setupJWTEnv(t)

	token, err := GenerateToken("testuser", "admin")
	require.NoError(t, err)

	claims, err := ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, "admin", claims.Role)
}

func TestValidateToken_ClaimsFields(t *testing.T) {
	setupJWTEnv(t)

	before := time.Now().Add(-time.Second)
	token, err := GenerateToken("alice", "super_admin")
	require.NoError(t, err)
	after := time.Now().Add(time.Second)

	claims, err := ValidateToken(token)
	require.NoError(t, err)

	assert.Equal(t, "alice", claims.Username)
	assert.Equal(t, "super_admin", claims.Role)

	// Check expiration is ~6 hours from now
	expTime := claims.ExpiresAt.Time
	assert.True(t, expTime.After(before.Add(5*time.Hour+59*time.Minute)),
		"expiration should be ~6 hours from now")
	assert.True(t, expTime.Before(after.Add(6*time.Hour+1*time.Minute)),
		"expiration should be ~6 hours from now")
}

func TestValidateToken_Expired(t *testing.T) {
	setupJWTEnv(t)

	// Set expiration to 0 hours
	t.Setenv("JWT_EXPIRATION_HOURS", "0")

	token, err := GenerateToken("testuser", "user")
	require.NoError(t, err)

	// Token with 0 hour expiration should be immediately expired
	// (IssuedAt and ExpiresAt are the same)
	// We need to wait a moment for it to be truly expired
	time.Sleep(time.Millisecond * 10)

	_, err = ValidateToken(token)
	assert.Error(t, err)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	setupJWTEnv(t)

	token, err := GenerateToken("testuser", "user")
	require.NoError(t, err)

	// Change the secret
	differentSecret := base64.StdEncoding.EncodeToString([]byte("different-secret-key-32bytes-xx!"))
	os.Setenv("JWT_SECRET", differentSecret)
	defer os.Setenv("JWT_SECRET", testSecret)

	_, err = ValidateToken(token)
	assert.Error(t, err)
}

func TestValidateToken_Malformed(t *testing.T) {
	setupJWTEnv(t)

	_, err := ValidateToken("not.a.valid.jwt.token")
	assert.Error(t, err)
}

func TestValidateToken_MissingSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	_, err := ValidateToken("some.token.here")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_SECRET not set")
}

func TestGenerateToken_UsesHS512(t *testing.T) {
	setupJWTEnv(t)

	tokenStr, err := GenerateToken("testuser", "user")
	require.NoError(t, err)

	// Parse without validation to check algorithm
	secret, _ := base64.StdEncoding.DecodeString(testSecret)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	require.NoError(t, err)
	assert.Equal(t, "HS512", token.Header["alg"])
}
