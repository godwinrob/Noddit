package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godwinrob/noddit/internal/models"
)

// VotePost handles voting on a post
func (h *Handler) VotePost(c *gin.Context) {
	var v models.Vote
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get user ID
	var userID int64
	err := h.db.QueryRow("SELECT id FROM users WHERE UPPER(username) = UPPER($1)", v.Username).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Determine vote value
	voteValue := 0
	if v.Vote == "upvote" {
		voteValue = 1
	} else if v.Vote == "downvote" {
		voteValue = -1
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vote type"})
		return
	}

	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer tx.Rollback()

	// Check if user has already voted
	var existingVote string
	err = tx.QueryRow("SELECT vote FROM post_votes WHERE post_id = $1 AND user_id = $2", v.PostID, userID).Scan(&existingVote)

	if err == nil {
		// User has already voted - update the vote
		_, err = tx.Exec("UPDATE post_votes SET vote = $1 WHERE post_id = $2 AND user_id = $3", v.Vote, v.PostID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vote"})
			return
		}

		// Adjust score based on vote change
		oldVoteValue := 0
		if existingVote == "upvote" {
			oldVoteValue = 1
		} else if existingVote == "downvote" {
			oldVoteValue = -1
		}
		scoreDiff := voteValue - oldVoteValue

		_, err = tx.Exec("UPDATE posts SET post_score = post_score + $1 WHERE post_id = $2", scoreDiff, v.PostID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update score"})
			return
		}
	} else {
		// New vote - insert and update score
		_, err = tx.Exec("INSERT INTO post_votes (post_id, user_id, vote) VALUES ($1, $2, $3)", v.PostID, userID, v.Vote)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record vote"})
			return
		}

		_, err = tx.Exec("UPDATE posts SET post_score = post_score + $1 WHERE post_id = $2", voteValue, v.PostID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update score"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote recorded"})
}

// GetVotes returns all votes for a post
func (h *Handler) GetVotes(c *gin.Context) {
	postID := c.Param("postId")

	rows, err := h.db.Query(`
		SELECT pv.post_id, pv.user_id, pv.vote, u.username
		FROM post_votes pv
		JOIN users u ON pv.user_id = u.id
		WHERE pv.post_id = $1
	`, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	votes := []models.Vote{}
	for rows.Next() {
		var v models.Vote
		if err := rows.Scan(&v.PostID, &v.UserID, &v.Vote, &v.Username); err != nil {
			continue
		}
		votes = append(votes, v)
	}

	c.JSON(http.StatusOK, votes)
}
