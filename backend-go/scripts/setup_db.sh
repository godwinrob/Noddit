#!/bin/bash

# Database setup script for Noddit
# Run this to create the database and apply migrations

set -e

# Load environment variables
if [ -f ../.env ]; then
    export $(cat ../.env | grep -v '^#' | xargs)
else
    echo "Warning: .env file not found, using defaults"
    DB_USER=${DB_USER:-postgres}
    DB_NAME=${DB_NAME:-userdb}
fi

echo "Setting up database: $DB_NAME"

# Create database if it doesn't exist
createdb -U $DB_USER $DB_NAME 2>/dev/null || echo "Database $DB_NAME already exists"

# Apply migrations
echo "Applying migrations..."
psql -U $DB_USER -d $DB_NAME -f ../migrations/001_initial_schema.up.sql

echo "✅ Database setup complete!"
echo ""
echo "To load sample data from the old database, run:"
echo "  psql -U $DB_USER -d $DB_NAME -f ../../backend/database/dbexport.pgsql"
