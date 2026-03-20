package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/godwinrob/noddit/pkg/auth"
)

// Context keys for storing user info
const (
	ContextKeyUsername = "username"
	ContextKeyRole     = "role"
	ContextKeyUserID   = "userId" // For future use if needed
)

// AuthMiddleware validates JWT tokens and requires authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			log.Printf("[Auth] Failed to extract token from %s %s: %v", c.Request.Method, c.Request.URL.Path, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}

		// Validate Clerk token
		claims, err := auth.ValidateClerkToken(token)
		if err != nil {
			// Log failed auth attempts for security monitoring
			log.Printf("[Auth] Invalid Clerk token attempt from IP %s on %s %s: %v",
				c.ClientIP(), c.Request.Method, c.Request.URL.Path, err)

			// Provide more specific error messages
			errorMsg := categorizeTokenError(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": errorMsg})
			c.Abort()
			return
		}

		// Extract username from Clerk claims (using Clerk user ID for now)
		username := claims.Subject

		// Get default role from environment variable
		defaultRole := os.Getenv("CLERK_DEFAULT_ROLE")
		if defaultRole == "" {
			defaultRole = "user" // Fallback if env var not set
		}

		// Store user info in context
		// Note: Clerk doesn't have "roles" by default, we'll need to use metadata or organizations
		c.Set(ContextKeyUsername, username)
		c.Set(ContextKeyRole, defaultRole) // Default role, can be enhanced with Clerk metadata

		c.Next()
	}
}

// AdminOnly middleware ensures user has admin or super_admin role
// Must be used AFTER AuthMiddleware
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextKeyRole)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Type-safe role checking
		roleStr, ok := role.(string)
		if !ok {
			log.Printf("[Auth] Invalid role type in context: %T", role)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
			c.Abort()
			return
		}

		if roleStr != "admin" && roleStr != "super_admin" {
			username, _ := c.Get(ContextKeyUsername)
			log.Printf("[Auth] Access denied: user %v (role: %s) attempted admin action on %s %s",
				username, roleStr, c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RoleRequired creates a middleware that checks for specific role(s)
func RoleRequired(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextKeyRole)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			log.Printf("[Auth] Invalid role type in context: %T", role)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
			c.Abort()
			return
		}

		// Check if user's role is in allowed list
		for _, allowed := range allowedRoles {
			if roleStr == allowed {
				c.Next()
				return
			}
		}

		username, _ := c.Get(ContextKeyUsername)
		log.Printf("[Auth] Access denied: user %v (role: %s) attempted action requiring roles %v on %s %s",
			username, roleStr, allowedRoles, c.Request.Method, c.Request.URL.Path)
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

// Helper functions for handlers to safely extract user info from context

// GetUsername returns the authenticated username from context
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get(ContextKeyUsername)
	if !exists {
		return "", false
	}
	usernameStr, ok := username.(string)
	return usernameStr, ok
}

// GetRole returns the authenticated user's role from context
func GetRole(c *gin.Context) (string, bool) {
	role, exists := c.Get(ContextKeyRole)
	if !exists {
		return "", false
	}
	roleStr, ok := role.(string)
	return roleStr, ok
}

// MustGetUsername returns the username or panics (use only after AuthMiddleware)
func MustGetUsername(c *gin.Context) string {
	username, ok := GetUsername(c)
	if !ok {
		panic("username not found in context - did you forget AuthMiddleware?")
	}
	return username
}

// IsAuthenticated checks if the current request has valid authentication
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get(ContextKeyUsername)
	return exists
}

// extractToken extracts the JWT token from the Authorization header
func extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	// Check Bearer prefix
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", errors.New("invalid authorization header format")
	}

	if parts[0] != "Bearer" {
		return "", errors.New("authorization header must use Bearer scheme")
	}

	if parts[1] == "" {
		return "", errors.New("token is empty")
	}

	return parts[1], nil
}

// categorizeTokenError provides user-friendly error messages based on Clerk errors
func categorizeTokenError(err error) string {
	// For Clerk errors, provide generic message (don't leak internal details)
	errStr := err.Error()
	if strings.Contains(errStr, "expired") {
		return "Token has expired"
	}
	return "Invalid or expired token"
}
