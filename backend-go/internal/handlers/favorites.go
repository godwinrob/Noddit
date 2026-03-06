package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godwinrob/noddit/internal/models"
)

// GetFavorites returns a user's favorite subnoddits
func (h *Handler) GetFavorites(c *gin.Context) {
	username := c.Param("username")

	rows, err := h.db.Query(`
		SELECT f.user_id, f.sn_id, s.sn_name
		FROM favorites f
		JOIN users u ON f.user_id = u.id
		LEFT JOIN subnoddits s ON f.sn_id = s.sn_id
		WHERE UPPER(u.username) = UPPER($1) AND f.sn_id IS NOT NULL
	`, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	favorites := []models.Favorites{}
	for rows.Next() {
		var f models.Favorites
		if err := rows.Scan(&f.UserID, &f.SubnodditID, &f.SubnodditName); err != nil {
			continue
		}
		favorites = append(favorites, f)
	}

	c.JSON(http.StatusOK, favorites)
}

// FavoritePost adds a post to favorites
func (h *Handler) FavoritePost(c *gin.Context) {
	var f models.Favorites
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get user ID
	var userID int64
	err := h.db.QueryRow("SELECT id FROM users WHERE UPPER(username) = UPPER($1)", f.Username).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Upsert favorite
	_, err = h.db.Exec(`
		INSERT INTO favorites (user_id, post_id)
		SELECT $1, $2
		WHERE NOT EXISTS (
			SELECT 1 FROM favorites WHERE user_id = $1 AND post_id = $2
		)
	`, userID, f.PostID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to favorite post"})
		return
	}

	// Return the favorite
	f.UserID = userID
	c.JSON(http.StatusOK, f)
}

// FavoriteSubnoddit adds a subnoddit to favorites
func (h *Handler) FavoriteSubnoddit(c *gin.Context) {
	var f models.Favorites
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get user ID
	var userID int64
	err := h.db.QueryRow("SELECT id FROM users WHERE UPPER(username) = UPPER($1)", f.Username).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Get subnoddit ID from name
	var subnodditID int64
	err = h.db.QueryRow("SELECT sn_id FROM subnoddits WHERE sn_name = $1", f.SubnodditName).Scan(&subnodditID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subnoddit not found"})
		return
	}

	// Upsert favorite
	_, err = h.db.Exec(`
		INSERT INTO favorites (user_id, sn_id)
		SELECT $1, $2
		WHERE NOT EXISTS (
			SELECT 1 FROM favorites WHERE user_id = $1 AND sn_id = $2
		)
	`, userID, subnodditID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to favorite subnoddit"})
		return
	}

	f.UserID = userID
	c.JSON(http.StatusOK, f)
}

// UnfavoritePost removes a post from favorites
func (h *Handler) UnfavoritePost(c *gin.Context) {
	postID := c.Param("postId")

	var f models.Favorites
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := h.db.Exec(`
		DELETE FROM favorites
		WHERE post_id = $1 AND user_id = $2
	`, postID, f.UserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfavorite post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post unfavorited"})
}

// UnfavoriteSubnoddit removes a subnoddit from favorites
func (h *Handler) UnfavoriteSubnoddit(c *gin.Context) {
	subnodditID := c.Param("subnodditId")

	// Get username from context (authenticated user)
	username, _ := c.Get("username")

	// Get user ID
	var userID int64
	err := h.db.QueryRow("SELECT id FROM users WHERE UPPER(username) = UPPER($1)", username).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	_, err = h.db.Exec(`
		DELETE FROM favorites
		WHERE sn_id = $1 AND user_id = $2
	`, subnodditID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfavorite subnoddit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subnoddit unfavorited"})
}
