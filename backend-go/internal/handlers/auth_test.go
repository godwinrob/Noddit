package handlers

import (
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/godwinrob/noddit/internal/models"
	"github.com/godwinrob/noddit/pkg/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var jwtTestSecret = base64.StdEncoding.EncodeToString([]byte("test-secret-for-jwt-testing-32b!"))

func TestLogin_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", jwtTestSecret)

	h, mock, err := setupMockDB()
	require.NoError(t, err)

	// Hash a real password so VerifyPassword works
	hash, salt, err := auth.HashPassword("testpass")
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "username", "password", "salt", "role"}).
		AddRow(1, "testuser", hash, salt, "user")
	mock.ExpectQuery("SELECT id, username, password, salt, role").
		WithArgs("testuser").
		WillReturnRows(rows)

	router := setupTestRouter()
	router.POST("/login", h.Login)

	w := performRequest(router, "POST", "/login", models.LoginRequest{
		Username: "testuser",
		Password: "testpass",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	// Response should be a JWT token string (quoted)
	assert.NotEmpty(t, w.Body.String())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_InvalidJSON(t *testing.T) {
	h, _, err := setupMockDB()
	require.NoError(t, err)

	router := setupTestRouter()
	router.POST("/login", h.Login)

	w := performRequest(router, "POST", "/login", "not-json")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_UserNotFound(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "username", "password", "salt", "role"})
	mock.ExpectQuery("SELECT id, username, password, salt, role").
		WithArgs("nonexistent").
		WillReturnRows(rows)

	router := setupTestRouter()
	router.POST("/login", h.Login)

	w := performRequest(router, "POST", "/login", models.LoginRequest{
		Username: "nonexistent",
		Password: "testpass",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_WrongPassword(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	hash, salt, err := auth.HashPassword("correctpass")
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "username", "password", "salt", "role"}).
		AddRow(1, "testuser", hash, salt, "user")
	mock.ExpectQuery("SELECT id, username, password, salt, role").
		WithArgs("testuser").
		WillReturnRows(rows)

	router := setupTestRouter()
	router.POST("/login", h.Login)

	w := performRequest(router, "POST", "/login", models.LoginRequest{
		Username: "testuser",
		Password: "wrongpass",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_DBError(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id, username, password, salt, role").
		WithArgs("testuser").
		WillReturnError(assert.AnError)

	router := setupTestRouter()
	router.POST("/login", h.Login)

	w := performRequest(router, "POST", "/login", models.LoginRequest{
		Username: "testuser",
		Password: "testpass",
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegister_Success(t *testing.T) {
	t.Setenv("DEFAULT_USER_ROLE", "user")
	t.Setenv("DEFAULT_FAVORITE_SUBNODDIT_ID", "1")

	h, mock, err := setupMockDB()
	require.NoError(t, err)

	// Username check
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	// Transaction
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("testuser", sqlmock.AnyArg(), sqlmock.AnyArg(), "user", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec("INSERT INTO favorites").
		WithArgs(int64(1), "1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := setupTestRouter()
	router.POST("/register", h.Register)

	w := performRequest(router, "POST", "/register", models.User{
		Username:        "testuser",
		Password:        "password123",
		ConfirmPassword: "password123",
	})

	assert.Equal(t, http.StatusOK, w.Code)

	var result models.RegistrationResult
	err = parseResponse(w, &result)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegister_InvalidJSON(t *testing.T) {
	h, _, err := setupMockDB()
	require.NoError(t, err)

	router := setupTestRouter()
	router.POST("/register", h.Register)

	w := performRequest(router, "POST", "/register", "bad-json")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_PasswordMismatch(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	// Username check still happens
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	router := setupTestRouter()
	router.POST("/register", h.Register)

	w := performRequest(router, "POST", "/register", models.User{
		Username:        "testuser",
		Password:        "password123",
		ConfirmPassword: "different123",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result models.RegistrationResult
	err = parseResponse(w, &result)
	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Contains(t, result.Errors, "Passwords do not match")
}

func TestRegister_PasswordTooShort(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	router := setupTestRouter()
	router.POST("/register", h.Register)

	w := performRequest(router, "POST", "/register", models.User{
		Username:        "testuser",
		Password:        "short",
		ConfirmPassword: "short",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result models.RegistrationResult
	err = parseResponse(w, &result)
	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Contains(t, result.Errors, "Password must be at least 8 characters")
}

func TestRegister_UsernameExists(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("existinguser").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	router := setupTestRouter()
	router.POST("/register", h.Register)

	w := performRequest(router, "POST", "/register", models.User{
		Username:        "existinguser",
		Password:        "password123",
		ConfirmPassword: "password123",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result models.RegistrationResult
	err = parseResponse(w, &result)
	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Contains(t, result.Errors, "Username already exists")
}

func TestRegister_MultipleErrors(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("existinguser").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	router := setupTestRouter()
	router.POST("/register", h.Register)

	w := performRequest(router, "POST", "/register", models.User{
		Username:        "existinguser",
		Password:        "short",
		ConfirmPassword: "diff",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result models.RegistrationResult
	err = parseResponse(w, &result)
	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.GreaterOrEqual(t, len(result.Errors), 2)
}

func TestRegister_DBError(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("testuser").
		WillReturnError(assert.AnError)

	router := setupTestRouter()
	router.POST("/register", h.Register)

	w := performRequest(router, "POST", "/register", models.User{
		Username:        "testuser",
		Password:        "password123",
		ConfirmPassword: "password123",
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
