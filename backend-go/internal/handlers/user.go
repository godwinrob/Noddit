package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godwinrob/noddit/internal/models"
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
