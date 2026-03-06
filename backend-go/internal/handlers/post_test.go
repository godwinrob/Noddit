package handlers

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/godwinrob/noddit/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var postColumns = []string{
	"post_id", "parent_post_id", "sn_id", "sn_name", "user_id", "username",
	"title", "body", "image_address", "created_date", "post_score", "top_level_id",
}

func addPostRow(rows *sqlmock.Rows, id int64, title string) *sqlmock.Rows {
	return rows.AddRow(
		id, nil, 1, "testsub", 1, "testuser",
		title, "Test body", nil, time.Now(), 5, nil,
	)
}

func TestGetAllPosts_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows(postColumns)
	addPostRow(rows, 1, "Post 1")
	addPostRow(rows, 2, "Post 2")

	mock.ExpectQuery("SELECT p.post_id").WillReturnRows(rows)

	router := setupTestRouter()
	router.GET("/posts", h.GetAllPosts)

	w := performRequest(router, "GET", "/posts", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var posts []models.Post
	err = parseResponse(w, &posts)
	require.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllPosts_Empty(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows(postColumns)
	mock.ExpectQuery("SELECT p.post_id").WillReturnRows(rows)

	router := setupTestRouter()
	router.GET("/posts", h.GetAllPosts)

	w := performRequest(router, "GET", "/posts", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var posts []models.Post
	err = parseResponse(w, &posts)
	require.NoError(t, err)
	assert.Len(t, posts, 0)
}

func TestGetAllPosts_DBError(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT p.post_id").WillReturnError(assert.AnError)

	router := setupTestRouter()
	router.GET("/posts", h.GetAllPosts)

	w := performRequest(router, "GET", "/posts", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetRecentPosts_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows(postColumns)
	addPostRow(rows, 1, "Recent Post")

	mock.ExpectQuery("SELECT p.post_id").WillReturnRows(rows)

	router := setupTestRouter()
	router.GET("/recent", h.GetRecentPosts)

	w := performRequest(router, "GET", "/recent", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var posts []models.Post
	err = parseResponse(w, &posts)
	require.NoError(t, err)
	assert.Len(t, posts, 1)
}

func TestGetPopularPosts_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows(postColumns)
	addPostRow(rows, 1, "Popular Post")

	mock.ExpectQuery("SELECT p.post_id").WillReturnRows(rows)

	router := setupTestRouter()
	router.GET("/popular", h.GetPopularPosts)

	w := performRequest(router, "GET", "/popular", nil)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetPopularPosts_DBError(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT p.post_id").WillReturnError(assert.AnError)

	router := setupTestRouter()
	router.GET("/popular", h.GetPopularPosts)

	w := performRequest(router, "GET", "/popular", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetPost_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows(postColumns)
	addPostRow(rows, 1, "Single Post")

	mock.ExpectQuery("SELECT p.post_id").
		WithArgs("1").
		WillReturnRows(rows)

	router := setupTestRouter()
	router.GET("/post/:postId", h.GetPost)

	w := performRequest(router, "GET", "/post/1", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var post models.Post
	err = parseResponse(w, &post)
	require.NoError(t, err)
	assert.Equal(t, "Single Post", post.Title)
}

func TestGetPost_NotFound(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT p.post_id").
		WithArgs("999").
		WillReturnError(sql.ErrNoRows)

	router := setupTestRouter()
	router.GET("/post/:postId", h.GetPost)

	w := performRequest(router, "GET", "/post/999", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreatePost_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	// Get user ID
	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Insert post
	mock.ExpectQuery("INSERT INTO posts").
		WithArgs(int64(1), int64(1), "Test Title", "Test Body", sql.NullString{}, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"post_id", "created_date"}).AddRow(1, time.Now()))

	router := setupTestRouter()
	router.POST("/post", h.CreatePost)

	w := performRequest(router, "POST", "/post", models.Post{
		SubnodditID: 1,
		Username:    "testuser",
		Title:       "Test Title",
		Body:        "Test Body",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePost_InvalidJSON(t *testing.T) {
	h, _, err := setupMockDB()
	require.NoError(t, err)

	router := setupTestRouter()
	router.POST("/post", h.CreatePost)

	w := performRequest(router, "POST", "/post", "bad")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePost_UserNotFound(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	router := setupTestRouter()
	router.POST("/post", h.CreatePost)

	w := performRequest(router, "POST", "/post", models.Post{
		SubnodditID: 1,
		Username:    "nonexistent",
		Title:       "Test",
		Body:        "Test",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateReply_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	// Get user ID
	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Insert reply
	mock.ExpectExec("INSERT INTO posts").
		WithArgs(int64(1), int64(1), int64(1), "Reply body", sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(2, 1))

	router := setupTestRouter()
	router.POST("/post/:postId/reply", h.CreateReply)

	// CreateReply uses ShouldBindJSON which requires title, username, body, subnodditId
	w := performRequest(router, "POST", "/post/1/reply", map[string]interface{}{
		"subnodditId": 1,
		"username":    "testuser",
		"title":       "reply",
		"body":        "Reply body",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeletePost_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM post_votes").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec("DELETE FROM posts").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	router := setupTestRouter()
	router.DELETE("/post/:postId", h.DeletePost)

	w := performRequest(router, "DELETE", "/post/1", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeletePost_DBError(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectBegin().WillReturnError(assert.AnError)

	router := setupTestRouter()
	router.DELETE("/post/:postId", h.DeletePost)

	w := performRequest(router, "DELETE", "/post/1", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetUserPosts_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows(postColumns)
	addPostRow(rows, 1, "User Post")

	mock.ExpectQuery("SELECT p.post_id").
		WithArgs("testuser").
		WillReturnRows(rows)

	router := setupTestRouter()
	router.GET("/user/:username/posts", h.GetUserPosts)

	w := performRequest(router, "GET", "/user/testuser/posts", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var posts []models.Post
	err = parseResponse(w, &posts)
	require.NoError(t, err)
	assert.Len(t, posts, 1)
}
