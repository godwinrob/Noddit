package handlers

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/godwinrob/noddit/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVotePost_NewUpvote(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	// Get user ID
	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectBegin()
	// Check existing vote - none found
	mock.ExpectQuery("SELECT vote FROM post_votes").
		WithArgs(int64(1), int64(1)).
		WillReturnError(sql.ErrNoRows)
	// Insert new vote
	mock.ExpectExec("INSERT INTO post_votes").
		WithArgs(int64(1), int64(1), "upvote").
		WillReturnResult(sqlmock.NewResult(1, 1))
	// Update score
	mock.ExpectExec("UPDATE posts SET post_score").
		WithArgs(1, int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// Response queries
	mock.ExpectQuery("SELECT post_score FROM posts").
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"post_score"}).AddRow(2))
	mock.ExpectQuery("SELECT vote FROM post_votes").
		WithArgs(int64(1), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"vote"}).AddRow("upvote"))

	router := setupTestRouter()
	router.POST("/vote", h.VotePost)

	w := performRequest(router, "POST", "/vote", models.Vote{
		PostID:   1,
		Username: "testuser",
		Vote:     "upvote",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVotePost_NewDownvote(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT vote FROM post_votes").
		WithArgs(int64(1), int64(1)).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO post_votes").
		WithArgs(int64(1), int64(1), "downvote").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE posts SET post_score").
		WithArgs(-1, int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// Response queries
	mock.ExpectQuery("SELECT post_score FROM posts").
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"post_score"}).AddRow(0))
	mock.ExpectQuery("SELECT vote FROM post_votes").
		WithArgs(int64(1), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"vote"}).AddRow("downvote"))

	router := setupTestRouter()
	router.POST("/vote", h.VotePost)

	w := performRequest(router, "POST", "/vote", models.Vote{
		PostID:   1,
		Username: "testuser",
		Vote:     "downvote",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVotePost_ChangeVote(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectBegin()
	// Existing vote found (downvote -> upvote = diff of 2)
	mock.ExpectQuery("SELECT vote FROM post_votes").
		WithArgs(int64(1), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"vote"}).AddRow("downvote"))
	// Update the vote
	mock.ExpectExec("UPDATE post_votes SET vote").
		WithArgs("upvote", int64(1), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	// Score diff: 1 - (-1) = 2
	mock.ExpectExec("UPDATE posts SET post_score").
		WithArgs(2, int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// Response queries
	mock.ExpectQuery("SELECT post_score FROM posts").
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"post_score"}).AddRow(3))
	mock.ExpectQuery("SELECT vote FROM post_votes").
		WithArgs(int64(1), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"vote"}).AddRow("upvote"))

	router := setupTestRouter()
	router.POST("/vote", h.VotePost)

	w := performRequest(router, "POST", "/vote", models.Vote{
		PostID:   1,
		Username: "testuser",
		Vote:     "upvote",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVotePost_ToggleOff(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectBegin()
	// Existing upvote found, same vote submitted — should toggle off
	mock.ExpectQuery("SELECT vote FROM post_votes").
		WithArgs(int64(1), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"vote"}).AddRow("upvote"))
	// Delete the vote
	mock.ExpectExec("DELETE FROM post_votes").
		WithArgs(int64(1), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	// Subtract the old vote value (1) from score
	mock.ExpectExec("UPDATE posts SET post_score").
		WithArgs(1, int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// Response queries — vote no longer exists
	mock.ExpectQuery("SELECT post_score FROM posts").
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"post_score"}).AddRow(1))
	mock.ExpectQuery("SELECT vote FROM post_votes").
		WithArgs(int64(1), int64(1)).
		WillReturnError(sql.ErrNoRows)

	router := setupTestRouter()
	router.POST("/vote", h.VotePost)

	w := performRequest(router, "POST", "/vote", models.Vote{
		PostID:   1,
		Username: "testuser",
		Vote:     "upvote",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVotePost_InvalidVoteType(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	router := setupTestRouter()
	router.POST("/vote", h.VotePost)

	w := performRequest(router, "POST", "/vote", models.Vote{
		PostID:   1,
		Username: "testuser",
		Vote:     "sidevote",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVotePost_UserNotFound(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	router := setupTestRouter()
	router.POST("/vote", h.VotePost)

	w := performRequest(router, "POST", "/vote", models.Vote{
		PostID:   1,
		Username: "nonexistent",
		Vote:     "upvote",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVotePost_InvalidJSON(t *testing.T) {
	h, _, err := setupMockDB()
	require.NoError(t, err)

	router := setupTestRouter()
	router.POST("/vote", h.VotePost)

	w := performRequest(router, "POST", "/vote", "bad")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetVotes_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"post_id", "user_id", "vote", "username"}).
		AddRow(1, 1, "upvote", "user1").
		AddRow(1, 2, "downvote", "user2")

	mock.ExpectQuery("SELECT pv.post_id").
		WithArgs("1").
		WillReturnRows(rows)

	router := setupTestRouter()
	router.GET("/votes/:postId", h.GetVotes)

	w := performRequest(router, "GET", "/votes/1", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var votes []models.Vote
	err = parseResponse(w, &votes)
	require.NoError(t, err)
	assert.Len(t, votes, 2)
}

func TestGetVotes_DBError(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT pv.post_id").
		WithArgs("1").
		WillReturnError(assert.AnError)

	router := setupTestRouter()
	router.GET("/votes/:postId", h.GetVotes)

	w := performRequest(router, "GET", "/votes/1", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
