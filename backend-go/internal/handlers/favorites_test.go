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

func TestGetFavorites_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"user_id", "sn_id", "sn_name"}).
		AddRow(1, 1, "golang").
		AddRow(1, 2, "react")

	mock.ExpectQuery("SELECT f.user_id").
		WithArgs("testuser").
		WillReturnRows(rows)

	router := setupTestRouter()
	router.GET("/favorites/:username", h.GetFavorites)

	w := performRequest(router, "GET", "/favorites/testuser", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var favorites []models.Favorites
	err = parseResponse(w, &favorites)
	require.NoError(t, err)
	assert.Len(t, favorites, 2)
}

func TestGetFavorites_Empty(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"user_id", "sn_id", "sn_name"})
	mock.ExpectQuery("SELECT f.user_id").
		WithArgs("testuser").
		WillReturnRows(rows)

	router := setupTestRouter()
	router.GET("/favorites/:username", h.GetFavorites)

	w := performRequest(router, "GET", "/favorites/testuser", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var favorites []models.Favorites
	err = parseResponse(w, &favorites)
	require.NoError(t, err)
	assert.Len(t, favorites, 0)
}

func TestGetFavorites_DBError(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT f.user_id").
		WithArgs("testuser").
		WillReturnError(assert.AnError)

	router := setupTestRouter()
	router.GET("/favorites/:username", h.GetFavorites)

	w := performRequest(router, "GET", "/favorites/testuser", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestFavoritePost_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec("INSERT INTO favorites").
		WithArgs(int64(1), sql.NullInt64{Int64: 42, Valid: true}).
		WillReturnResult(sqlmock.NewResult(1, 1))

	router := setupTestRouter()
	router.POST("/favorite/post", h.FavoritePost)

	w := performRequest(router, "POST", "/favorite/post", models.Favorites{
		Username: "testuser",
		PostID:   sql.NullInt64{Int64: 42, Valid: true},
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFavoriteSubnoddit_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery("SELECT sn_id FROM subnoddits").
		WithArgs("golang").
		WillReturnRows(sqlmock.NewRows([]string{"sn_id"}).AddRow(5))

	mock.ExpectExec("INSERT INTO favorites").
		WithArgs(int64(1), int64(5)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	router := setupTestRouter()
	router.POST("/favorite/subnoddit", h.FavoriteSubnoddit)

	w := performRequest(router, "POST", "/favorite/subnoddit", models.Favorites{
		Username:      "testuser",
		SubnodditName: "golang",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnfavoritePost_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectExec("DELETE FROM favorites").
		WithArgs("42", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	router := setupTestRouter()
	router.DELETE("/favorite/post/:postId", h.UnfavoritePost)

	w := performRequest(router, "DELETE", "/favorite/post/42", models.Favorites{
		UserID:   1,
		Username: "testuser",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnfavoriteSubnoddit_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	// Get user ID from username in request body
	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec("DELETE FROM favorites").
		WithArgs("5", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	router := setupTestRouter()
	router.DELETE("/favorite/subnoddit/:subnodditId", h.UnfavoriteSubnoddit)

	w := performRequest(router, "DELETE", "/favorite/subnoddit/5", models.Favorites{
		Username: "testuser",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFavoritePost_UserNotFound(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	router := setupTestRouter()
	router.POST("/favorite/post", h.FavoritePost)

	w := performRequest(router, "POST", "/favorite/post", models.Favorites{
		Username: "nonexistent",
		PostID:   sql.NullInt64{Int64: 42, Valid: true},
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
