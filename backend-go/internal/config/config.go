package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	// Server
	Port        string
	FrontendURL string

	// Database
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	DBMaxOpenConns int
	DBMaxIdleConns int

	// Clerk
	ClerkSecretKey   string
	ClerkDefaultRole string
}

// Load reads configuration from environment variables and validates them
func Load() (*Config, error) {
	cfg := &Config{
		// Server
		Port:        getEnv("PORT", "8080"),
		FrontendURL: os.Getenv("FRONTEND_URL"),

		// Database
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		DBMaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 5),

		// Clerk
		ClerkSecretKey:   os.Getenv("CLERK_SECRET_KEY"),
		ClerkDefaultRole: getEnv("CLERK_DEFAULT_ROLE", "user"),
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Log configuration (mask sensitive values)
	cfg.LogConfig()

	return cfg, nil
}

// Validate checks that all required configuration is present
func (c *Config) Validate() error {
	required := map[string]string{
		"FRONTEND_URL": c.FrontendURL,
		"DB_HOST":      c.DBHost,
		"DB_PORT":      c.DBPort,
		"DB_USER":      c.DBUser,
		"DB_PASSWORD":  c.DBPassword,
		"DB_NAME":      c.DBName,
	}

	var missing []string
	for key, value := range required {
		if value == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}

	// Validate connection pool values
	if c.DBMaxOpenConns < 1 {
		return fmt.Errorf("DB_MAX_OPEN_CONNS must be at least 1")
	}
	if c.DBMaxIdleConns < 0 {
		return fmt.Errorf("DB_MAX_IDLE_CONNS must be 0 or greater")
	}
	if c.DBMaxIdleConns > c.DBMaxOpenConns {
		return fmt.Errorf("DB_MAX_IDLE_CONNS cannot be greater than DB_MAX_OPEN_CONNS")
	}

	return nil
}

// LogConfig logs the current configuration (masking sensitive values)
func (c *Config) LogConfig() {
	log.Println("=== Application Configuration ===")
	log.Printf("Server Port: %s", c.Port)
	log.Printf("Frontend URL: %s", c.FrontendURL)
	log.Printf("Database Host: %s:%s", c.DBHost, c.DBPort)
	log.Printf("Database Name: %s", c.DBName)
	log.Printf("Database User: %s", c.DBUser)
	log.Printf("Database SSL Mode: %s", c.DBSSLMode)
	log.Printf("DB Max Open Connections: %d", c.DBMaxOpenConns)
	log.Printf("DB Max Idle Connections: %d", c.DBMaxIdleConns)
	log.Printf("Clerk Secret Key: %s", maskSecret(c.ClerkSecretKey))
	log.Printf("Clerk Default Role: %s", c.ClerkDefaultRole)
	log.Println("=================================")
}

// Helper functions

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s, using default: %d", key, defaultValue)
		return defaultValue
	}

	return value
}

func maskSecret(secret string) string {
	if secret == "" {
		return "<not set - keyless mode>"
	}
	if len(secret) < 8 {
		return "***"
	}
	return secret[:4] + "..." + secret[len(secret)-4:]
}
