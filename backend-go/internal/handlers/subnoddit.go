package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/godwinrob/noddit/internal/models"
)

// GetAllSubnoddits returns all subnoddits
func (h *Handler) GetAllSubnoddits(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT sn_id, sn_name, sn_description
		FROM subnoddits
		ORDER BY sn_name
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	subnoddits := []models.Subnoddit{}
	for rows.Next() {
		var s models.Subnoddit
		if err := rows.Scan(&s.SubnodditID, &s.SubnodditName, &s.SubnodditDescription); err != nil {
			continue
		}
		subnoddits = append(subnoddits, s)
	}

	c.JSON(http.StatusOK, subnoddits)
}

// GetActiveSubnoddits returns the 5 most active subnoddits
func (h *Handler) GetActiveSubnoddits(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT s.sn_id, s.sn_name, s.sn_description, MAX(p.post_id) as post_id
		FROM subnoddits s
		JOIN posts p ON s.sn_id = p.sn_id
		GROUP BY s.sn_id, s.sn_name, s.sn_description
		ORDER BY MAX(p.post_id) DESC
		LIMIT 5
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	subnoddits := []models.Subnoddit{}
	for rows.Next() {
		var s models.Subnoddit
		if err := rows.Scan(&s.SubnodditID, &s.SubnodditName, &s.SubnodditDescription, &s.PostID); err != nil {
			continue
		}
		subnoddits = append(subnoddits, s)
	}

	c.JSON(http.StatusOK, subnoddits)
}

// GetSubnodditByName returns a specific subnoddit by name
func (h *Handler) GetSubnodditByName(c *gin.Context) {
	name := c.Param("name")

	var s models.Subnoddit
	err := h.db.QueryRow(`
		SELECT sn_id, sn_name, sn_description
		FROM subnoddits
		WHERE sn_name = $1
	`, name).Scan(&s.SubnodditID, &s.SubnodditName, &s.SubnodditDescription)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subnoddit not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, s)
}

// SearchSubnoddits searches for subnoddits by term
func (h *Handler) SearchSubnoddits(c *gin.Context) {
	term := c.Param("term")
	searchTerm := "%" + term + "%"

	rows, err := h.db.Query(`
		SELECT sn_id, sn_name, sn_description
		FROM subnoddits
		WHERE sn_name ILIKE $1 OR sn_description ILIKE $1
		ORDER BY sn_name
	`, searchTerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	subnoddits := []models.Subnoddit{}
	for rows.Next() {
		var s models.Subnoddit
		if err := rows.Scan(&s.SubnodditID, &s.SubnodditName, &s.SubnodditDescription); err != nil {
			continue
		}
		subnoddits = append(subnoddits, s)
	}

	c.JSON(http.StatusOK, subnoddits)
}

// CreateSubnoddit creates a new subnoddit
func (h *Handler) CreateSubnoddit(c *gin.Context) {
	var s models.Subnoddit
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Replace spaces with underscores
	s.SubnodditName = strings.ReplaceAll(s.SubnodditName, " ", "_")

	// Get user ID from username
	var userID int64
	err := h.db.QueryRow("SELECT id FROM users WHERE UPPER(username) = UPPER($1)", s.Username).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer tx.Rollback()

	// Insert subnoddit
	var subnodditID int64
	err = tx.QueryRow(`
		INSERT INTO subnoddits (sn_name, sn_description)
		VALUES ($1, $2)
		RETURNING sn_id
	`, s.SubnodditName, s.SubnodditDescription).Scan(&subnodditID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subnoddit"})
		return
	}

	// Make creator a moderator
	_, err = tx.Exec(`
		INSERT INTO mod (sn_id, user_id)
		VALUES ($1, $2)
	`, subnodditID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add moderator"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete creation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subnoddit created"})
}

// UpdateSubnoddit updates a subnoddit's description
func (h *Handler) UpdateSubnoddit(c *gin.Context) {
	var s models.Subnoddit
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := h.db.Exec(`
		UPDATE subnoddits
		SET sn_description = $1
		WHERE sn_id = $2
	`, s.SubnodditDescription, s.SubnodditID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subnoddit"})
		return
	}

	c.JSON(http.StatusOK, s)
}

// DeleteSubnoddit deletes a subnoddit
func (h *Handler) DeleteSubnoddit(c *gin.Context) {
	name := c.Param("name")

	_, err := h.db.Exec("DELETE FROM subnoddits WHERE sn_name = $1", name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subnoddit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subnoddit deleted"})
}

// GetModerators returns moderators for a subnoddit
func (h *Handler) GetModerators(c *gin.Context) {
	name := c.Param("name")

	rows, err := h.db.Query(`
		SELECT m.sn_id, m.user_id, u.username
		FROM mod m
		JOIN subnoddits s ON m.sn_id = s.sn_id
		JOIN users u ON m.user_id = u.id
		WHERE s.sn_name = $1
	`, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	mods := []models.Moderator{}
	for rows.Next() {
		var m models.Moderator
		if err := rows.Scan(&m.SubnodditID, &m.UserID, &m.Username); err != nil {
			continue
		}
		mods = append(mods, m)
	}

	c.JSON(http.StatusOK, mods)
}
