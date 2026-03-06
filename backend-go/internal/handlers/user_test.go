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

func TestGetUser_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "username", "role", "avatar_address", "first_name", "last_name", "email_address", "join_date"}).
		AddRow(1, "testuser", "user", nil, nil, nil, nil, time.Now())

	mock.ExpectQuery("SELECT id, username, role").
		WithArgs("testuser").
		WillReturnRows(rows)

	router := setupTestRouter()
	router.GET("/user/:username", h.GetUser)

	w := performRequest(router, "GET", "/user/testuser", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var user models.User
	err = parseResponse(w, &user)
	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
}

func TestGetUser_NotFound(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id, username, role").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	router := setupTestRouter()
	router.GET("/user/:username", h.GetUser)

	w := performRequest(router, "GET", "/user/nonexistent", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUser_DBError(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id, username, role").
		WithArgs("testuser").
		WillReturnError(assert.AnError)

	router := setupTestRouter()
	router.GET("/user/:username", h.GetUser)

	w := performRequest(router, "GET", "/user/testuser", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateUsername_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectExec("UPDATE users").
		WithArgs("newname", "oldname").
		WillReturnResult(sqlmock.NewResult(0, 1))

	router := setupTestRouter()
	router.PUT("/user/:username/username", h.UpdateUsername)

	w := performRequest(router, "PUT", "/user/oldname/username", models.User{
		Username:    "oldname",
		Password:    "dummy",
		NewUsername: "newname",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateName_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "username", "role", "avatar_address", "first_name", "last_name", "email_address", "join_date"}).
		AddRow(1, "testuser", "user", nil, "John", "Doe", nil, time.Now())

	mock.ExpectQuery("UPDATE users").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "testuser").
		WillReturnRows(rows)

	router := setupTestRouter()
	router.PUT("/user/:username/name", h.UpdateName)

	// User model has binding:"required" on username and password
	w := performRequest(router, "PUT", "/user/testuser/name", map[string]interface{}{
		"username":  "testuser",
		"password":  "dummy",
		"firstName": map[string]interface{}{"String": "John", "Valid": true},
		"lastName":  map[string]interface{}{"String": "Doe", "Valid": true},
	})

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateAvatar_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "username", "role", "avatar_address", "first_name", "last_name", "email_address", "join_date"}).
		AddRow(1, "testuser", "user", "https://example.com/avatar.jpg", nil, nil, nil, time.Now())

	mock.ExpectQuery("UPDATE users").
		WithArgs(sqlmock.AnyArg(), "testuser").
		WillReturnRows(rows)

	router := setupTestRouter()
	router.PUT("/user/:username/avatar", h.UpdateAvatar)

	w := performRequest(router, "PUT", "/user/testuser/avatar", map[string]interface{}{
		"username":      "testuser",
		"password":      "dummy",
		"avatarAddress": map[string]interface{}{"String": "https://example.com/avatar.jpg", "Valid": true},
	})

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteUser_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectExec("DELETE FROM users").
		WithArgs("testuser").
		WillReturnResult(sqlmock.NewResult(0, 1))

	router := setupTestRouter()
	router.DELETE("/user/:username", h.DeleteUser)

	w := performRequest(router, "DELETE", "/user/testuser", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser_DBError(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectExec("DELETE FROM users").
		WithArgs("testuser").
		WillReturnError(assert.AnError)

	router := setupTestRouter()
	router.DELETE("/user/:username", h.DeleteUser)

	w := performRequest(router, "DELETE", "/user/testuser", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
