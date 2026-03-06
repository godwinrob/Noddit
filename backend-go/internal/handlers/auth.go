package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godwinrob/noddit/internal/models"
	"github.com/godwinrob/noddit/pkg/auth"
)

// Login handles user login and returns JWT token
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get user from database
	var user models.User
	query := `SELECT id, username, password, salt, role
	          FROM users
	          WHERE UPPER(username) = UPPER($1)`

	err := h.db.QueryRow(query, req.Username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Salt,
		&user.Role,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Verify password
	valid, err := auth.VerifyPassword(req.Password, user.Password, user.Salt)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, token)
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.RegistrationResult{
			Success: false,
			Errors:  []string{"Invalid request"},
		})
		return
	}

	// Validation
	errors := []string{}

	if user.Password != user.ConfirmPassword {
		errors = append(errors, "Passwords do not match")
	}

	if len(user.Password) < 8 {
		errors = append(errors, "Password must be at least 8 characters")
	}

	user.Username = strings.ToLower(user.Username)

	// Check if username exists
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE UPPER(username) = UPPER($1))", user.Username).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.RegistrationResult{
			Success: false,
			Errors:  []string{"Database error"},
		})
		return
	}

	if exists {
		errors = append(errors, "Username already exists")
	}

	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, models.RegistrationResult{
			Success: false,
			Errors:  errors,
		})
		return
	}

	// Hash password
	hashedPassword, salt, err := auth.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.RegistrationResult{
			Success: false,
			Errors:  []string{"Failed to hash password"},
		})
		return
	}

	// Set default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	// Insert user
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.RegistrationResult{
			Success: false,
			Errors:  []string{"Database error"},
		})
		return
	}
	defer tx.Rollback()

	var userID int64
	err = tx.QueryRow(`
		INSERT INTO users (username, password, salt, role, join_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		user.Username, hashedPassword, salt, user.Role, time.Now(),
	).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.RegistrationResult{
			Success: false,
			Errors:  []string{"Failed to create user"},
		})
		return
	}

	// Add default favorite (subnoddit with id 1 - "Cats")
	_, err = tx.Exec(`
		INSERT INTO favorites (user_id, sn_id)
		SELECT $1, 1
		WHERE NOT EXISTS (
			SELECT 1 FROM favorites WHERE user_id = $1 AND sn_id = 1
		)`,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.RegistrationResult{
			Success: false,
			Errors:  []string{"Failed to add default favorite"},
		})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.RegistrationResult{
			Success: false,
			Errors:  []string{"Failed to complete registration"},
		})
		return
	}

	c.JSON(http.StatusOK, models.RegistrationResult{
		Success: true,
		Errors:  []string{},
	})
}
