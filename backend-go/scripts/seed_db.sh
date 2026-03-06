#!/bin/bash

# Seed database with sample data
# Run this AFTER setup_db.sh to populate with example posts

set -e

# Load environment variables
if [ -f ../.env ]; then
    export $(cat ../.env | grep -v '^#' | xargs)
else
    echo "Warning: .env file not found, using defaults"
    DB_USER=${DB_USER:-postgres}
    DB_NAME=${DB_NAME:-userdb}
fi

echo "Seeding database: $DB_NAME with sample data..."

# Apply seed migration
psql -U $DB_USER -d $DB_NAME -f ../migrations/002_seed_data.up.sql

echo "✅ Database seeded successfully!"
echo ""
echo "Sample users created (username = password):"
echo "  - rgodwin (super_admin)"
echo "  - csamad (super_admin)"
echo "  - eknutson (super_admin)"
echo "  - jminihan (super_admin)"
echo "  - test (user)"
echo "  - asd1 (user)"
echo "  - david (user)"
echo ""
echo "Sample communities:"
echo "  - Cats"
echo "  - Dogs"
echo "  - Harold"
echo "  - Gardening"
echo "  - star_wars"
echo ""
echo "🎉 Ready to explore!"
