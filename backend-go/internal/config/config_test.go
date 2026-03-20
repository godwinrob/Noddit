package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("FRONTEND_URL", "http://localhost:8081")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "testuser")
	t.Setenv("DB_PASSWORD", "testpass")
	t.Setenv("DB_NAME", "testdb")
}

func TestLoad_WithAllVars(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("PORT", "9090")
	t.Setenv("DB_SSLMODE", "require")
	t.Setenv("DB_MAX_OPEN_CONNS", "50")
	t.Setenv("DB_MAX_IDLE_CONNS", "10")

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "http://localhost:8081", cfg.FrontendURL)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "testuser", cfg.DBUser)
	assert.Equal(t, "testpass", cfg.DBPassword)
	assert.Equal(t, "testdb", cfg.DBName)
	assert.Equal(t, "require", cfg.DBSSLMode)
	assert.Equal(t, 50, cfg.DBMaxOpenConns)
	assert.Equal(t, 10, cfg.DBMaxIdleConns)
}

func TestLoad_DefaultValues(t *testing.T) {
	setRequiredEnv(t)
	// Don't set optional vars - should use defaults
	t.Setenv("PORT", "")
	t.Setenv("DB_SSLMODE", "")
	t.Setenv("DB_MAX_OPEN_CONNS", "")
	t.Setenv("DB_MAX_IDLE_CONNS", "")

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "disable", cfg.DBSSLMode)
	assert.Equal(t, 25, cfg.DBMaxOpenConns)
	assert.Equal(t, 5, cfg.DBMaxIdleConns)
	assert.Equal(t, "user", cfg.ClerkDefaultRole)
}

func TestLoad_MissingRequiredVars(t *testing.T) {
	// Clear all required vars
	t.Setenv("FRONTEND_URL", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "")

	_, err := Load()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required environment variables")
}

func TestValidate_PoolMaxOpenConnsZero(t *testing.T) {
	cfg := &Config{
		FrontendURL:    "http://localhost",
		DBHost:         "localhost",
		DBPort:         "5432",
		DBUser:         "user",
		DBPassword:     "pass",
		DBName:         "db",
		DBMaxOpenConns: 0,
		DBMaxIdleConns: 0,
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DB_MAX_OPEN_CONNS must be at least 1")
}

func TestValidate_PoolIdleGreaterThanOpen(t *testing.T) {
	cfg := &Config{
		FrontendURL:    "http://localhost",
		DBHost:         "localhost",
		DBPort:         "5432",
		DBUser:         "user",
		DBPassword:     "pass",
		DBName:         "db",
		DBMaxOpenConns: 5,
		DBMaxIdleConns: 10,
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DB_MAX_IDLE_CONNS cannot be greater than DB_MAX_OPEN_CONNS")
}

func TestValidate_NegativeIdleConns(t *testing.T) {
	cfg := &Config{
		FrontendURL:    "http://localhost",
		DBHost:         "localhost",
		DBPort:         "5432",
		DBUser:         "user",
		DBPassword:     "pass",
		DBName:         "db",
		DBMaxOpenConns: 5,
		DBMaxIdleConns: -1,
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DB_MAX_IDLE_CONNS must be 0 or greater")
}

func TestValidate_ValidPoolSettings(t *testing.T) {
	cfg := &Config{
		FrontendURL:    "http://localhost",
		DBHost:         "localhost",
		DBPort:         "5432",
		DBUser:         "user",
		DBPassword:     "pass",
		DBName:         "db",
		DBMaxOpenConns: 25,
		DBMaxIdleConns: 5,
	}

	err := cfg.Validate()
	assert.NoError(t, err)
}
