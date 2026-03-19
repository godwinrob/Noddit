package handlers

import (
	"database/sql"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/godwinrob/noddit/internal/models"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_]{2,19}$`)
	reservedNames = map[string]bool{
		"admin":     true,
		"moderator": true,
		"noddit":    true,
		"system":    true,
		"deleted":   true,
	}
)

// GetUser returns a user profile
func (h *Handler) GetUser(c *gin.Context) {
	username := c.Param("username")

	var u models.User
	err := h.db.QueryRow(`
		SELECT id, username, role, avatar_address, first_name, last_name, email_address, join_date
		FROM users
		WHERE UPPER(username) = UPPER($1)
	`, username).Scan(&u.ID, &u.Username, &u.Role, &u.AvatarAddress, &u.FirstName, &u.LastName, &u.EmailAddress, &u.JoinDate)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, u)
}

// UpdateEmail updates a user's email
func (h *Handler) UpdateEmail(c *gin.Context) {
	username := c.Param("username")

	var email string
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := h.db.Exec(`
		UPDATE users
		SET email_address = $1
		WHERE UPPER(username) = UPPER($2)
	`, email, username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email updated"})
}

// UpdateUsername updates a user's username
func (h *Handler) UpdateUsername(c *gin.Context) {
	username := c.Param("username")

	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := h.db.Exec(`
		UPDATE users
		SET username = $1
		WHERE UPPER(username) = UPPER($2)
	`, u.NewUsername, username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Username updated"})
}

// UpdateName updates a user's first and last name
func (h *Handler) UpdateName(c *gin.Context) {
	username := c.Param("username")

	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.db.QueryRow(`
		UPDATE users
		SET first_name = $1, last_name = $2
		WHERE UPPER(username) = UPPER($3)
		RETURNING id, username, role, avatar_address, first_name, last_name, email_address, join_date
	`, u.FirstName, u.LastName, username).Scan(
		&u.ID, &u.Username, &u.Role, &u.AvatarAddress,
		&u.FirstName, &u.LastName, &u.EmailAddress, &u.JoinDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update name"})
		return
	}

	c.JSON(http.StatusOK, u)
}

// UpdateAvatar updates a user's avatar
func (h *Handler) UpdateAvatar(c *gin.Context) {
	username := c.Param("username")

	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.db.QueryRow(`
		UPDATE users
		SET avatar_address = $1
		WHERE UPPER(username) = UPPER($2)
		RETURNING id, username, role, avatar_address, first_name, last_name, email_address, join_date
	`, u.AvatarAddress, username).Scan(
		&u.ID, &u.Username, &u.Role, &u.AvatarAddress,
		&u.FirstName, &u.LastName, &u.EmailAddress, &u.JoinDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar"})
		return
	}

	c.JSON(http.StatusOK, u)
}

// SyncUser ensures a Clerk-authenticated user exists in the database.
// Called by the frontend after sign-in as a post-auth hook.
//
// Modes:
//   - checkOnly=true: look up by email, return {isNew:true} or {isNew:false, username, userId}
//   - checkOnly=false (default): create the user with the given username
func (h *Handler) SyncUser(c *gin.Context) {
	var req struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		CheckOnly bool   `json:"checkOnly"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	// Look up existing user by email (stable Clerk identity)
	var userID int64
	var username string
	err := h.db.QueryRow(
		"SELECT id, username FROM users WHERE LOWER(email_address) = LOWER($1)", req.Email,
	).Scan(&userID, &username)

	if err == nil {
		// User exists
		c.JSON(http.StatusOK, gin.H{"isNew": false, "username": username, "userId": userID})
		return
	}
	if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// User does not exist
	if req.CheckOnly {
		c.JSON(http.StatusOK, gin.H{"isNew": true})
		return
	}

	// Create mode — username is required
	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	if !usernameRegex.MatchString(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be 3-20 characters, alphanumeric or underscore, and cannot start with underscore"})
		return
	}
	if reservedNames[strings.ToLower(req.Username)] {
		c.JSON(http.StatusConflict, gin.H{"error": "Username is reserved"})
		return
	}

	// Insert user
	err = h.db.QueryRow(`
		INSERT INTO users (username, password, salt, role, email_address, join_date)
		VALUES (LOWER($1), 'clerk-managed', 'clerk-managed', 'user', $2, NOW())
		RETURNING id`, req.Username, req.Email).Scan(&userID)
	if err != nil {
		// Check if it's a unique constraint violation (username taken)
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Username is already taken"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isNew": false, "username": strings.ToLower(req.Username), "userId": userID})
}

// CheckUsernameAvailable checks if a username is available for registration.
func (h *Handler) CheckUsernameAvailable(c *gin.Context) {
	username := c.Param("username")

	// Validate format
	if !usernameRegex.MatchString(username) {
		c.JSON(http.StatusOK, gin.H{
			"available": false,
			"reason":    "Must be 3-20 characters, letters/numbers/underscores only, cannot start with underscore",
		})
		return
	}

	// Check reserved names
	if reservedNames[strings.ToLower(username)] {
		c.JSON(http.StatusOK, gin.H{
			"available": false,
			"reason":    "This username is reserved",
		})
		return
	}

	// Check DB uniqueness (case-insensitive)
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1))", username).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if exists {
		c.JSON(http.StatusOK, gin.H{
			"available": false,
			"reason":    "Username is already taken",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"available": true})
}

// DeleteUser deletes a user (admin only)
func (h *Handler) DeleteUser(c *gin.Context) {
	username := c.Param("username")

	_, err := h.db.Exec("DELETE FROM users WHERE UPPER(username) = UPPER($1)", username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
