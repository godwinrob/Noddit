package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestContext(authHeader string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	if authHeader != "" {
		c.Request.Header.Set("Authorization", authHeader)
	}
	return c, w
}

// --- extractToken tests ---

func TestExtractToken_ValidBearer(t *testing.T) {
	c, _ := setupTestContext("Bearer mytoken123")
	token, err := extractToken(c)
	assert.NoError(t, err)
	assert.Equal(t, "mytoken123", token)
}

func TestExtractToken_MissingHeader(t *testing.T) {
	c, _ := setupTestContext("")
	_, err := extractToken(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing authorization header")
}

func TestExtractToken_NoPrefix(t *testing.T) {
	c, _ := setupTestContext("mytoken123")
	_, err := extractToken(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid authorization header format")
}

func TestExtractToken_WrongScheme(t *testing.T) {
	c, _ := setupTestContext("Basic mytoken123")
	_, err := extractToken(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Bearer scheme")
}

// --- AdminOnly tests ---

func TestAdminOnly_AdminPasses(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Set(ContextKeyRole, "admin")
	c.Set(ContextKeyUsername, "adminuser")

	handler := AdminOnly()
	handler(c)

	assert.False(t, c.IsAborted())
}

func TestAdminOnly_SuperAdminPasses(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Set(ContextKeyRole, "super_admin")
	c.Set(ContextKeyUsername, "superadmin")

	handler := AdminOnly()
	handler(c)

	assert.False(t, c.IsAborted())
}

func TestAdminOnly_UserGetsForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Set(ContextKeyRole, "user")
	c.Set(ContextKeyUsername, "regularuser")

	handler := AdminOnly()
	handler(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAdminOnly_NoRoleGetsUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	// Don't set any role

	handler := AdminOnly()
	handler(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// --- RoleRequired tests ---

func TestRoleRequired_AllowedRole(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Set(ContextKeyRole, "moderator")
	c.Set(ContextKeyUsername, "moduser")

	handler := RoleRequired("admin", "moderator")
	handler(c)

	assert.False(t, c.IsAborted())
}

func TestRoleRequired_DeniedRole(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Set(ContextKeyRole, "user")
	c.Set(ContextKeyUsername, "regularuser")

	handler := RoleRequired("admin", "moderator")
	handler(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRoleRequired_NoRole(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	handler := RoleRequired("admin")
	handler(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// --- Helper function tests ---

func TestGetUsername_Exists(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ContextKeyUsername, "testuser")

	username, ok := GetUsername(c)
	assert.True(t, ok)
	assert.Equal(t, "testuser", username)
}

func TestGetUsername_NotExists(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	_, ok := GetUsername(c)
	assert.False(t, ok)
}

func TestIsAuthenticated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	assert.False(t, IsAuthenticated(c))

	c.Set(ContextKeyUsername, "testuser")
	assert.True(t, IsAuthenticated(c))
}
