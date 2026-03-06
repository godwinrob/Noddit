package handlers

import (
	"database/sql"
)

// Handler holds the database connection
type Handler struct {
	db *sql.DB
}

// NewHandler creates a new handler with database connection
func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}
