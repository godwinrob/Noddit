package models

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostMarshalJSON_NullFieldsOmitted(t *testing.T) {
	p := Post{
		PostID:      1,
		SubnodditID: 1,
		UserID:      1,
		Username:    "testuser",
		Title:       "Test Post",
		Body:        "Test body",
		CreatedDate: time.Now(),
		PostScore:   5,
		// ParentPostID, ImageAddress, TopLevelID are null
	}

	data, err := json.Marshal(p)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Nil(t, result["imageAddress"], "null imageAddress should be omitted")
	assert.Nil(t, result["parentPostId"], "null parentPostId should be omitted")
	assert.Nil(t, result["topLevelId"], "null topLevelId should be omitted")
}

func TestPostMarshalJSON_ValidFieldsIncluded(t *testing.T) {
	imgAddr := "https://example.com/image.jpg"
	p := Post{
		PostID:       1,
		ParentPostID: sql.NullInt64{Int64: 42, Valid: true},
		SubnodditID:  1,
		UserID:       1,
		Username:     "testuser",
		Title:        "Test Post",
		Body:         "Test body",
		ImageAddress: sql.NullString{String: imgAddr, Valid: true},
		CreatedDate:  time.Now(),
		PostScore:    5,
		TopLevelID:   sql.NullInt64{Int64: 10, Valid: true},
	}

	data, err := json.Marshal(p)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, imgAddr, result["imageAddress"])
	assert.Equal(t, float64(42), result["parentPostId"])
	assert.Equal(t, float64(10), result["topLevelId"])
}

func TestUserMarshalJSON_NullFieldsOmitted(t *testing.T) {
	u := User{
		ID:       1,
		Username: "testuser",
		Role:     "user",
		// AvatarAddress, FirstName, LastName, EmailAddress are null
	}

	data, err := json.Marshal(u)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Nil(t, result["avatarAddress"], "null avatarAddress should be omitted")
	assert.Nil(t, result["firstName"], "null firstName should be omitted")
	assert.Nil(t, result["lastName"], "null lastName should be omitted")
	assert.Nil(t, result["emailAddress"], "null emailAddress should be omitted")
}

func TestUserMarshalJSON_ValidFieldsIncluded(t *testing.T) {
	u := User{
		ID:            1,
		Username:      "testuser",
		Role:          "user",
		AvatarAddress: sql.NullString{String: "https://example.com/avatar.jpg", Valid: true},
		FirstName:     sql.NullString{String: "John", Valid: true},
		LastName:      sql.NullString{String: "Doe", Valid: true},
		EmailAddress:  sql.NullString{String: "john@example.com", Valid: true},
	}

	data, err := json.Marshal(u)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "https://example.com/avatar.jpg", result["avatarAddress"])
	assert.Equal(t, "John", result["firstName"])
	assert.Equal(t, "Doe", result["lastName"])
	assert.Equal(t, "john@example.com", result["emailAddress"])
}
