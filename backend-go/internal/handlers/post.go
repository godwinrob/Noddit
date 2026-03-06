package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godwinrob/noddit/internal/models"
)

// GetAllPosts returns all posts
func (h *Handler) GetAllPosts(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT p.post_id, p.parent_post_id, p.sn_id, s.sn_name, p.user_id, u.username,
		       p.title, p.body, p.image_address, p.created_date, p.post_score, p.top_level_id
		FROM posts p
		JOIN subnoddits s ON p.sn_id = s.sn_id
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_date DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.PostID, &p.ParentPostID, &p.SubnodditID, &p.SubnodditName,
			&p.UserID, &p.Username, &p.Title, &p.Body, &p.ImageAddress,
			&p.CreatedDate, &p.PostScore, &p.TopLevelID); err != nil {
			continue
		}
		posts = append(posts, p)
	}

	c.JSON(http.StatusOK, posts)
}

// GetRecentPosts returns 5 most recent posts
func (h *Handler) GetRecentPosts(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT p.post_id, p.parent_post_id, p.sn_id, s.sn_name, p.user_id, u.username,
		       p.title, p.body, p.image_address, p.created_date, p.post_score, p.top_level_id
		FROM posts p
		JOIN subnoddits s ON p.sn_id = s.sn_id
		JOIN users u ON p.user_id = u.id
		ORDER BY p.post_id DESC
		LIMIT 5
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.PostID, &p.ParentPostID, &p.SubnodditID, &p.SubnodditName,
			&p.UserID, &p.Username, &p.Title, &p.Body, &p.ImageAddress,
			&p.CreatedDate, &p.PostScore, &p.TopLevelID); err != nil {
			continue
		}
		posts = append(posts, p)
	}

	c.JSON(http.StatusOK, posts)
}

// GetPopularPosts returns 10 popular posts from last 24 hours
func (h *Handler) GetPopularPosts(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT p.post_id, p.parent_post_id, p.sn_id, s.sn_name, p.user_id, u.username,
		       p.title, p.body, p.image_address, p.created_date, p.post_score, p.top_level_id
		FROM posts p
		JOIN subnoddits s ON p.sn_id = s.sn_id
		JOIN users u ON p.user_id = u.id
		WHERE p.created_date >= NOW() - INTERVAL '24 hours'
		ORDER BY p.post_score DESC, p.created_date DESC
		LIMIT 10
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.PostID, &p.ParentPostID, &p.SubnodditID, &p.SubnodditName,
			&p.UserID, &p.Username, &p.Title, &p.Body, &p.ImageAddress,
			&p.CreatedDate, &p.PostScore, &p.TopLevelID); err != nil {
			continue
		}
		posts = append(posts, p)
	}

	c.JSON(http.StatusOK, posts)
}

// GetPostsBySubnoddit returns all posts for a subnoddit
func (h *Handler) GetPostsBySubnoddit(c *gin.Context) {
	subnodditName := c.Param("subnodditName")

	rows, err := h.db.Query(`
		SELECT p.post_id, p.parent_post_id, p.sn_id, s.sn_name, p.user_id, u.username,
		       p.title, p.body, p.image_address, p.created_date, p.post_score, p.top_level_id
		FROM posts p
		JOIN subnoddits s ON p.sn_id = s.sn_id
		JOIN users u ON p.user_id = u.id
		WHERE s.sn_name = $1
		ORDER BY p.created_date DESC
	`, subnodditName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.PostID, &p.ParentPostID, &p.SubnodditID, &p.SubnodditName,
			&p.UserID, &p.Username, &p.Title, &p.Body, &p.ImageAddress,
			&p.CreatedDate, &p.PostScore, &p.TopLevelID); err != nil {
			continue
		}
		posts = append(posts, p)
	}

	c.JSON(http.StatusOK, posts)
}

// GetPopularPostsBySubnoddit returns posts for a subnoddit sorted by score
func (h *Handler) GetPopularPostsBySubnoddit(c *gin.Context) {
	subnodditName := c.Param("subnodditName")

	rows, err := h.db.Query(`
		SELECT p.post_id, p.parent_post_id, p.sn_id, s.sn_name, p.user_id, u.username,
		       p.title, p.body, p.image_address, p.created_date, p.post_score, p.top_level_id
		FROM posts p
		JOIN subnoddits s ON p.sn_id = s.sn_id
		JOIN users u ON p.user_id = u.id
		WHERE s.sn_name = $1
		ORDER BY p.post_score DESC
	`, subnodditName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.PostID, &p.ParentPostID, &p.SubnodditID, &p.SubnodditName,
			&p.UserID, &p.Username, &p.Title, &p.Body, &p.ImageAddress,
			&p.CreatedDate, &p.PostScore, &p.TopLevelID); err != nil {
			continue
		}
		posts = append(posts, p)
	}

	c.JSON(http.StatusOK, posts)
}

// GetPost returns a single post by ID
func (h *Handler) GetPost(c *gin.Context) {
	postID := c.Param("postId")

	var p models.Post
	err := h.db.QueryRow(`
		SELECT p.post_id, p.parent_post_id, p.sn_id, s.sn_name, p.user_id, u.username,
		       p.title, p.body, p.image_address, p.created_date, p.post_score, p.top_level_id
		FROM posts p
		JOIN subnoddits s ON p.sn_id = s.sn_id
		JOIN users u ON p.user_id = u.id
		WHERE p.post_id = $1
	`, postID).Scan(&p.PostID, &p.ParentPostID, &p.SubnodditID, &p.SubnodditName,
		&p.UserID, &p.Username, &p.Title, &p.Body, &p.ImageAddress,
		&p.CreatedDate, &p.PostScore, &p.TopLevelID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, p)
}

// GetReplies returns all replies to a post
func (h *Handler) GetReplies(c *gin.Context) {
	postID := c.Param("postId")

	rows, err := h.db.Query(`
		SELECT p.post_id, p.parent_post_id, p.sn_id, s.sn_name, p.user_id, u.username,
		       p.title, p.body, p.image_address, p.created_date, p.post_score, p.top_level_id
		FROM posts p
		JOIN subnoddits s ON p.sn_id = s.sn_id
		JOIN users u ON p.user_id = u.id
		WHERE p.top_level_id = $1
		ORDER BY p.post_score DESC
	`, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.PostID, &p.ParentPostID, &p.SubnodditID, &p.SubnodditName,
			&p.UserID, &p.Username, &p.Title, &p.Body, &p.ImageAddress,
			&p.CreatedDate, &p.PostScore, &p.TopLevelID); err != nil {
			continue
		}
		posts = append(posts, p)
	}

	c.JSON(http.StatusOK, posts)
}

// CreatePost creates a new post
func (h *Handler) CreatePost(c *gin.Context) {
	var p models.Post
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get user ID
	var userID int64
	err := h.db.QueryRow("SELECT id FROM users WHERE UPPER(username) = UPPER($1)", p.Username).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Insert post
	err = h.db.QueryRow(`
		INSERT INTO posts (sn_id, user_id, title, body, image_address, created_date, post_score)
		VALUES ($1, $2, $3, $4, $5, $6, 1)
		RETURNING post_id, created_date
	`, p.SubnodditID, userID, p.Title, p.Body, p.ImageAddress, time.Now()).Scan(&p.PostID, &p.CreatedDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusOK, p)
}

// UpdatePost updates a post's body and image
func (h *Handler) UpdatePost(c *gin.Context) {
	postID := c.Param("postId")

	var p models.Post
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := h.db.Exec(`
		UPDATE posts
		SET body = $1, image_address = $2
		WHERE post_id = $3
	`, p.Body, p.ImageAddress, postID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	// Get updated post
	err = h.db.QueryRow(`
		SELECT post_id, parent_post_id, sn_id, user_id, title, body, image_address, created_date, post_score, top_level_id
		FROM posts WHERE post_id = $1
	`, postID).Scan(&p.PostID, &p.ParentPostID, &p.SubnodditID, &p.UserID,
		&p.Title, &p.Body, &p.ImageAddress, &p.CreatedDate, &p.PostScore, &p.TopLevelID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated post"})
		return
	}

	c.JSON(http.StatusOK, p)
}

// CreateReply creates a reply to a post
func (h *Handler) CreateReply(c *gin.Context) {
	postID := c.Param("postId")

	var p models.Post
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get user ID
	var userID int64
	err := h.db.QueryRow("SELECT id FROM users WHERE UPPER(username) = UPPER($1)", p.Username).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Convert postID string to int64
	parentPostID, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Determine if this is a post-reply or comment-reply
	var topLevelID int64
	if p.ParentPostID.Valid {
		// This is a comment-reply
		topLevelID = p.TopLevelID.Int64
	} else {
		// This is a post-reply
		topLevelID = parentPostID
	}

	// Insert reply
	_, err = h.db.Exec(`
		INSERT INTO posts (parent_post_id, sn_id, user_id, title, body, created_date, post_score, top_level_id)
		VALUES ($1, $2, $3, '', $4, $5, 1, $6)
	`, parentPostID, p.SubnodditID, userID, p.Body, time.Now(), topLevelID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reply"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reply created"})
}

// DeletePost deletes a post
func (h *Handler) DeletePost(c *gin.Context) {
	postID := c.Param("postId")

	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer tx.Rollback()

	// Delete votes first
	_, err = tx.Exec("DELETE FROM post_votes WHERE post_id = $1", postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete votes"})
		return
	}

	// Delete post
	_, err = tx.Exec("DELETE FROM posts WHERE post_id = $1", postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete deletion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

// GetUserPosts returns all posts by a user
func (h *Handler) GetUserPosts(c *gin.Context) {
	username := c.Param("username")

	rows, err := h.db.Query(`
		SELECT p.post_id, p.parent_post_id, p.sn_id, s.sn_name, p.user_id, u.username,
		       p.title, p.body, p.image_address, p.created_date, p.post_score, p.top_level_id
		FROM posts p
		JOIN subnoddits s ON p.sn_id = s.sn_id
		JOIN users u ON p.user_id = u.id
		WHERE UPPER(u.username) = UPPER($1)
		ORDER BY p.created_date DESC
	`, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.PostID, &p.ParentPostID, &p.SubnodditID, &p.SubnodditName,
			&p.UserID, &p.Username, &p.Title, &p.Body, &p.ImageAddress,
			&p.CreatedDate, &p.PostScore, &p.TopLevelID); err != nil {
			continue
		}
		posts = append(posts, p)
	}

	c.JSON(http.StatusOK, posts)
}
