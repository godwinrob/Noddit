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

func TestGetAllSubnoddits_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"sn_id", "sn_name", "sn_description"}).
		AddRow(1, "golang", "Go programming").
		AddRow(2, "react", "React development")

	mock.ExpectQuery("SELECT sn_id, sn_name, sn_description").WillReturnRows(rows)

	router := setupTestRouter()
	router.GET("/subnoddits", h.GetAllSubnoddits)

	w := performRequest(router, "GET", "/subnoddits", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var subnoddits []models.Subnoddit
	err = parseResponse(w, &subnoddits)
	require.NoError(t, err)
	assert.Len(t, subnoddits, 2)
}

func TestGetAllSubnoddits_Empty(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"sn_id", "sn_name", "sn_description"})
	mock.ExpectQuery("SELECT sn_id, sn_name, sn_description").WillReturnRows(rows)

	router := setupTestRouter()
	router.GET("/subnoddits", h.GetAllSubnoddits)

	w := performRequest(router, "GET", "/subnoddits", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var subnoddits []models.Subnoddit
	err = parseResponse(w, &subnoddits)
	require.NoError(t, err)
	assert.Len(t, subnoddits, 0)
}

func TestGetAllSubnoddits_DBError(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT sn_id, sn_name, sn_description").WillReturnError(assert.AnError)

	router := setupTestRouter()
	router.GET("/subnoddits", h.GetAllSubnoddits)

	w := performRequest(router, "GET", "/subnoddits", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetSubnodditByName_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"sn_id", "sn_name", "sn_description"}).
		AddRow(1, "golang", "Go programming")

	mock.ExpectQuery("SELECT sn_id, sn_name, sn_description").
		WithArgs("golang").
		WillReturnRows(rows)

	router := setupTestRouter()
	router.GET("/subnoddit/:name", h.GetSubnodditByName)

	w := performRequest(router, "GET", "/subnoddit/golang", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var sn models.Subnoddit
	err = parseResponse(w, &sn)
	require.NoError(t, err)
	assert.Equal(t, "golang", sn.SubnodditName)
}

func TestGetSubnodditByName_NotFound(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT sn_id, sn_name, sn_description").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	router := setupTestRouter()
	router.GET("/subnoddit/:name", h.GetSubnodditByName)

	w := performRequest(router, "GET", "/subnoddit/nonexistent", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSearchSubnoddits_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"sn_id", "sn_name", "sn_description"}).
		AddRow(1, "golang", "Go programming language")

	mock.ExpectQuery("SELECT sn_id, sn_name, sn_description").
		WithArgs("%go%").
		WillReturnRows(rows)

	router := setupTestRouter()
	router.GET("/search/:term", h.SearchSubnoddits)

	w := performRequest(router, "GET", "/search/go", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var subnoddits []models.Subnoddit
	err = parseResponse(w, &subnoddits)
	require.NoError(t, err)
	assert.Len(t, subnoddits, 1)
}

func TestCreateSubnoddit_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	// Get user ID
	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Transaction
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO subnoddits").
		WithArgs("my_subnoddit", "A test community").
		WillReturnRows(sqlmock.NewRows([]string{"sn_id"}).AddRow(1))
	mock.ExpectExec("INSERT INTO mod").
		WithArgs(int64(1), int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := setupTestRouter()
	router.POST("/subnoddit", h.CreateSubnoddit)

	w := performRequest(router, "POST", "/subnoddit", models.Subnoddit{
		SubnodditName:        "my subnoddit",
		SubnodditDescription: "A test community",
		Username:             "testuser",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateSubnoddit_SpaceToUnderscore(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectBegin()
	// The name should have spaces replaced with underscores
	mock.ExpectQuery("INSERT INTO subnoddits").
		WithArgs("my_cool_sub", "desc").
		WillReturnRows(sqlmock.NewRows([]string{"sn_id"}).AddRow(1))
	mock.ExpectExec("INSERT INTO mod").
		WithArgs(int64(1), int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := setupTestRouter()
	router.POST("/subnoddit", h.CreateSubnoddit)

	w := performRequest(router, "POST", "/subnoddit", models.Subnoddit{
		SubnodditName:        "my cool sub",
		SubnodditDescription: "desc",
		Username:             "testuser",
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateSubnoddit_UserNotFound(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	router := setupTestRouter()
	router.POST("/subnoddit", h.CreateSubnoddit)

	w := performRequest(router, "POST", "/subnoddit", models.Subnoddit{
		SubnodditName:        "test",
		SubnodditDescription: "desc",
		Username:             "nonexistent",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSubnoddit_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectExec("UPDATE subnoddits").
		WithArgs("Updated description", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	router := setupTestRouter()
	router.PUT("/subnoddit", h.UpdateSubnoddit)

	w := performRequest(router, "PUT", "/subnoddit", models.Subnoddit{
		SubnodditID:          1,
		SubnodditName:        "test",
		SubnodditDescription: "Updated description",
	})

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteSubnoddit_Success(t *testing.T) {
	h, mock, err := setupMockDB()
	require.NoError(t, err)

	mock.ExpectExec("DELETE FROM subnoddits").
		WithArgs("test").
		WillReturnResult(sqlmock.NewResult(0, 1))

	router := setupTestRouter()
	router.DELETE("/subnoddit/:name", h.DeleteSubnoddit)

	w := performRequest(router, "DELETE", "/subnoddit/test", nil)

	assert.Equal(t, http.StatusOK, w.Code)
}
