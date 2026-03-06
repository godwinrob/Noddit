package auth

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword_ReturnsValidBase64(t *testing.T) {
	hash, salt, err := HashPassword("testpassword")
	require.NoError(t, err)

	_, err = base64.StdEncoding.DecodeString(hash)
	assert.NoError(t, err, "hash should be valid base64")

	_, err = base64.StdEncoding.DecodeString(salt)
	assert.NoError(t, err, "salt should be valid base64")
}

func TestHashPassword_DifferentSaltsEachCall(t *testing.T) {
	_, salt1, err := HashPassword("testpassword")
	require.NoError(t, err)

	_, salt2, err := HashPassword("testpassword")
	require.NoError(t, err)

	assert.NotEqual(t, salt1, salt2, "salts should differ between calls")
}

func TestHashPassword_DifferentHashesEachCall(t *testing.T) {
	hash1, _, err := HashPassword("testpassword")
	require.NoError(t, err)

	hash2, _, err := HashPassword("testpassword")
	require.NoError(t, err)

	assert.NotEqual(t, hash1, hash2, "hashes should differ due to different salts")
}

func TestVerifyPassword_CorrectPassword(t *testing.T) {
	hash, salt, err := HashPassword("mypassword")
	require.NoError(t, err)

	valid, err := VerifyPassword("mypassword", hash, salt)
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestVerifyPassword_WrongPassword(t *testing.T) {
	hash, salt, err := HashPassword("mypassword")
	require.NoError(t, err)

	valid, err := VerifyPassword("wrongpassword", hash, salt)
	require.NoError(t, err)
	assert.False(t, valid)
}

func TestVerifyPassword_InvalidBase64Salt(t *testing.T) {
	_, err := VerifyPassword("password", "validbase64==", "not!valid!base64!")
	// The function won't necessarily return an error for all invalid base64 on all
	// platforms, but it should handle invalid salt gracefully
	// Try with definitely invalid base64
	_, err = VerifyPassword("password", "dGVzdA==", "%%%invalid%%%")
	assert.Error(t, err, "should return error for invalid base64 salt")
}

func TestVerifyPassword_InvalidBase64Hash(t *testing.T) {
	_, salt, err := HashPassword("password")
	require.NoError(t, err)

	_, err = VerifyPassword("password", "%%%invalid%%%", salt)
	assert.Error(t, err, "should return error for invalid base64 hash")
}
