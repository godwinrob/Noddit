package handlers

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
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
	router.PUT("/user/:username/username", func(c *gin.Context) {
		c.Set("username", "oldname") // Mock auth context
		h.UpdateUsername(c)
	})

	w := performRequest(router, "PUT", "/user/oldname/username", models.UpdateUsernameRequest{
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
	router.PUT("/user/:username/name", func(c *gin.Context) {
		c.Set("username", "testuser") // Mock auth context
		h.UpdateName(c)
	})

	w := performRequest(router, "PUT", "/user/testuser/name", models.UpdateNameRequest{
		FirstName: sql.NullString{String: "John", Valid: true},
		LastName:  sql.NullString{String: "Doe", Valid: true},
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
	router.PUT("/user/:username/avatar", func(c *gin.Context) {
		c.Set("username", "testuser") // Mock auth context
		h.UpdateAvatar(c)
	})

	w := performRequest(router, "PUT", "/user/testuser/avatar", models.UpdateAvatarRequest{
		AvatarAddress: sql.NullString{String: "https://example.com/avatar.jpg", Valid: true},
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

// --- SyncUser tests ---

func TestSyncUser_CheckOnly_ExistingUser(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "testuser")
	mock.ExpectQuery("SELECT id, username FROM users WHERE LOWER").
		WithArgs("test@example.com").
		WillReturnRows(rows)

	router := setupTestRouter()
	router.POST("/api/user/sync", h.SyncUser)

	w := performRequest(router, "POST", "/api/user/sync", map[string]interface{}{
		"email":     "test@example.com",
		"checkOnly": true,
	})

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = parseResponse(w, &resp)
	require.NoError(t, err)
	assert.Equal(t, false, resp["isNew"])
	assert.Equal(t, "testuser", resp["username"])
}

func TestSyncUser_CheckOnly_NewUser(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id, username FROM users WHERE LOWER").
		WithArgs("new@example.com").
		WillReturnError(sql.ErrNoRows)

	router := setupTestRouter()
	router.POST("/api/user/sync", h.SyncUser)

	w := performRequest(router, "POST", "/api/user/sync", map[string]interface{}{
		"email":     "new@example.com",
		"checkOnly": true,
	})

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = parseResponse(w, &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["isNew"])
}

func TestSyncUser_Create_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id, username FROM users WHERE LOWER").
		WithArgs("new@example.com").
		WillReturnError(sql.ErrNoRows)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(42)
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("cooluser", "new@example.com").
		WillReturnRows(rows)

	router := setupTestRouter()
	router.POST("/api/user/sync", h.SyncUser)

	w := performRequest(router, "POST", "/api/user/sync", map[string]interface{}{
		"email":    "new@example.com",
		"username": "cooluser",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = parseResponse(w, &resp)
	require.NoError(t, err)
	assert.Equal(t, false, resp["isNew"])
	assert.Equal(t, "cooluser", resp["username"])
}

func TestSyncUser_Create_MissingUsername(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id, username FROM users WHERE LOWER").
		WithArgs("new@example.com").
		WillReturnError(sql.ErrNoRows)

	router := setupTestRouter()
	router.POST("/api/user/sync", h.SyncUser)

	w := performRequest(router, "POST", "/api/user/sync", map[string]interface{}{
		"email": "new@example.com",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSyncUser_Create_InvalidUsername(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id, username FROM users WHERE LOWER").
		WithArgs("new@example.com").
		WillReturnError(sql.ErrNoRows)

	router := setupTestRouter()
	router.POST("/api/user/sync", h.SyncUser)

	w := performRequest(router, "POST", "/api/user/sync", map[string]interface{}{
		"email":    "new@example.com",
		"username": "_bad",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSyncUser_Create_ReservedUsername(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id, username FROM users WHERE LOWER").
		WithArgs("new@example.com").
		WillReturnError(sql.ErrNoRows)

	router := setupTestRouter()
	router.POST("/api/user/sync", h.SyncUser)

	w := performRequest(router, "POST", "/api/user/sync", map[string]interface{}{
		"email":    "new@example.com",
		"username": "admin",
	})

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestSyncUser_MissingEmail(t *testing.T) {
	h, _, err := setupMockDB()
	require.NoError(t, err)

	router := setupTestRouter()
	router.POST("/api/user/sync", h.SyncUser)

	w := performRequest(router, "POST", "/api/user/sync", map[string]interface{}{
		"username": "test",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- CheckUsernameAvailable tests ---

func TestCheckUsernameAvailable_Available(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("newuser").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	router := setupTestRouter()
	router.GET("/api/public/user/available/:username", h.CheckUsernameAvailable)

	w := performRequest(router, "GET", "/api/public/user/available/newuser", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = parseResponse(w, &resp)
	require.NoError(t, err)
	assert.Equal(t, true, resp["available"])
}

func TestCheckUsernameAvailable_Taken(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("existinguser").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	router := setupTestRouter()
	router.GET("/api/public/user/available/:username", h.CheckUsernameAvailable)

	w := performRequest(router, "GET", "/api/public/user/available/existinguser", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = parseResponse(w, &resp)
	require.NoError(t, err)
	assert.Equal(t, false, resp["available"])
}

func TestCheckUsernameAvailable_InvalidFormat(t *testing.T) {
	h, _, err := setupMockDB()
	require.NoError(t, err)

	router := setupTestRouter()
	router.GET("/api/public/user/available/:username", h.CheckUsernameAvailable)

	w := performRequest(router, "GET", "/api/public/user/available/_bad", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = parseResponse(w, &resp)
	require.NoError(t, err)
	assert.Equal(t, false, resp["available"])
}

func TestCheckUsernameAvailable_Reserved(t *testing.T) {
	h, _, err := setupMockDB()
	require.NoError(t, err)

	router := setupTestRouter()
	router.GET("/api/public/user/available/:username", h.CheckUsernameAvailable)

	w := performRequest(router, "GET", "/api/public/user/available/admin", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = parseResponse(w, &resp)
	require.NoError(t, err)
	assert.Equal(t, false, resp["available"])
	assert.Equal(t, "This username is reserved", resp["reason"])
}

// --- Ownership Check Tests ---

func TestUpdateEmail_OwnershipDenied(t *testing.T) {
	h, _, err := setupMockDB()
	require.NoError(t, err)

	router := setupTestRouter()
	router.PUT("/user/:username/email", func(c *gin.Context) {
		c.Set("username", "alice") // Alice trying to update bob's email
		h.UpdateEmail(c)
	})

	w := performRequest(router, "PUT", "/user/bob/email", "new@email.com")

	assert.Equal(t, http.StatusForbidden, w.Code)
	var resp map[string]interface{}
	err = parseResponse(w, &resp)
	require.NoError(t, err)
	assert.Equal(t, "You can only update your own account", resp["error"])
}

func TestUpdateUsername_OwnershipDenied(t *testing.T) {
	h, _, err := setupMockDB()
	require.NoError(t, err)

	router := setupTestRouter()
	router.PUT("/user/:username/username", func(c *gin.Context) {
		c.Set("username", "alice") // Alice trying to update bob's username
		h.UpdateUsername(c)
	})

	w := performRequest(router, "PUT", "/user/bob/username", models.UpdateUsernameRequest{
		NewUsername: "bobby",
	})

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUpdateName_OwnershipDenied(t *testing.T) {
	h, _, err := setupMockDB()
	require.NoError(t, err)

	router := setupTestRouter()
	router.PUT("/user/:username/name", func(c *gin.Context) {
		c.Set("username", "alice")
		h.UpdateName(c)
	})

	w := performRequest(router, "PUT", "/user/bob/name", models.UpdateNameRequest{
		FirstName: sql.NullString{String: "Bob", Valid: true},
		LastName:  sql.NullString{String: "Smith", Valid: true},
	})

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUpdateAvatar_OwnershipDenied(t *testing.T) {
	h, _, err := setupMockDB()
	require.NoError(t, err)

	router := setupTestRouter()
	router.PUT("/user/:username/avatar", func(c *gin.Context) {
		c.Set("username", "alice")
		h.UpdateAvatar(c)
	})

	w := performRequest(router, "PUT", "/user/bob/avatar", models.UpdateAvatarRequest{
		AvatarAddress: sql.NullString{String: "https://example.com/avatar.jpg", Valid: true},
	})

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUpdateEmail_CaseInsensitiveOwnership(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectExec("UPDATE users").
		WithArgs("new@email.com", "TestUser").
		WillReturnResult(sqlmock.NewResult(0, 1))

	router := setupTestRouter()
	router.PUT("/user/:username/email", func(c *gin.Context) {
		c.Set("username", "testuser") // lowercase in token
		h.UpdateEmail(c)
	})

	w := performRequest(router, "PUT", "/user/TestUser/email", "new@email.com")

	assert.Equal(t, http.StatusOK, w.Code)
}
