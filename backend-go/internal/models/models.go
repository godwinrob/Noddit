package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// User represents a user in the system
type User struct {
	ID            int64          `json:"id"`
	Username      string         `json:"username" binding:"required,max=20"`
	Password      string         `json:"password,omitempty" binding:"required,max=100"`
	ConfirmPassword string       `json:"confirmPassword,omitempty"`
	Salt          string         `json:"-"` // Never send to client
	Role          string         `json:"role"`
	AvatarAddress sql.NullString `json:"avatarAddress"`
	FirstName     sql.NullString `json:"firstName"`
	LastName      sql.NullString `json:"lastName"`
	EmailAddress  sql.NullString `json:"emailAddress"`
	JoinDate      sql.NullTime   `json:"joinDate"`
	NewUsername   string         `json:"newUsername,omitempty"`
}

// Post represents a post or comment
type Post struct {
	PostID         int64          `json:"postId"`
	ParentPostID   sql.NullInt64  `json:"parentPostId"`
	SubnodditID    int64          `json:"subnodditId" binding:"required"`
	SubnodditName  string         `json:"subnodditName,omitempty"`
	UserID         int64          `json:"userId"`
	Username       string         `json:"username" binding:"required,max=20"`
	Title          string         `json:"title" binding:"required,min=1,max=300"`
	Body           string         `json:"body" binding:"required,min=1,max=10000"`
	ImageAddress   sql.NullString `json:"imageAddress"`
	CreatedDate    time.Time      `json:"createdDate"`
	PostScore      int64          `json:"postScore"`
	TopLevelID     sql.NullInt64  `json:"topLevelId"`
}

// Subnoddit represents a subnoddit (like a subreddit)
type Subnoddit struct {
	SubnodditID          int64  `json:"subnodditId"`
	SubnodditName        string `json:"subnodditName" binding:"required,min=3,max=50"`
	SubnodditDescription string `json:"subnodditDescription" binding:"required,min=1,max=500"`
	Username             string `json:"username,omitempty"`
	PostID               sql.NullInt64 `json:"postId,omitempty"`
}

// Vote represents a vote on a post
type Vote struct {
	UserID   int64  `json:"userId"`
	PostID   int64  `json:"postId" binding:"required"`
	Vote     string `json:"vote" binding:"required"` // "upvote" or "downvote"
	Username string `json:"username" binding:"required"`
}

// Favorites represents a user's favorite post or subnoddit
type Favorites struct {
	UserID         int64         `json:"userId"`
	SubnodditID    sql.NullInt64 `json:"subnodditId"`
	PostID         sql.NullInt64 `json:"postId"`
	SubnodditName  string        `json:"subnodditName,omitempty"`
	Username       string        `json:"username" binding:"required"`
}

// Moderator represents a subnoddit moderator
type Moderator struct {
	SubnodditID int64  `json:"subnodditId"`
	UserID      int64  `json:"userId"`
	Username    string `json:"username"`
}

// RegistrationResult represents the result of a registration attempt
type RegistrationResult struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse represents the JWT token response
type TokenResponse struct {
	Token string `json:"token"`
}

// UpdateUsernameRequest represents a request to update username
type UpdateUsernameRequest struct {
	NewUsername string `json:"newUsername" binding:"required,max=20"`
}

// UpdateNameRequest represents a request to update name
type UpdateNameRequest struct {
	FirstName sql.NullString `json:"firstName"`
	LastName  sql.NullString `json:"lastName"`
}

// UpdateAvatarRequest represents a request to update avatar
type UpdateAvatarRequest struct {
	AvatarAddress sql.NullString `json:"avatarAddress"`
}

// MarshalJSON customizes JSON marshaling for Post to handle nullable fields
func (p Post) MarshalJSON() ([]byte, error) {
	type Alias Post
	return json.Marshal(&struct {
		*Alias
		ImageAddress   *string `json:"imageAddress,omitempty"`
		ParentPostID   *int64  `json:"parentPostId,omitempty"`
		TopLevelID     *int64  `json:"topLevelId,omitempty"`
	}{
		Alias:        (*Alias)(&p),
		ImageAddress: nullStringToPtr(p.ImageAddress),
		ParentPostID: nullInt64ToPtr(p.ParentPostID),
		TopLevelID:   nullInt64ToPtr(p.TopLevelID),
	})
}

// MarshalJSON customizes JSON marshaling for User to handle nullable fields
func (u User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		*Alias
		AvatarAddress *string `json:"avatarAddress,omitempty"`
		FirstName     *string `json:"firstName,omitempty"`
		LastName      *string `json:"lastName,omitempty"`
		EmailAddress  *string `json:"emailAddress,omitempty"`
	}{
		Alias:         (*Alias)(&u),
		AvatarAddress: nullStringToPtr(u.AvatarAddress),
		FirstName:     nullStringToPtr(u.FirstName),
		LastName:      nullStringToPtr(u.LastName),
		EmailAddress:  nullStringToPtr(u.EmailAddress),
	})
}

// Helper functions to convert sql.Null* to pointers
func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullInt64ToPtr(ni sql.NullInt64) *int64 {
	if ni.Valid {
		return &ni.Int64
	}
	return nil
}
