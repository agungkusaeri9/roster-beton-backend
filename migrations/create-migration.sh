#!/bin/bash

# Script to create new migration files
# Usage: ./create-migration.sh <migration_name>
# Example: ./create-migration.sh add_email_to_users

if [ -z "$1" ]; then
    echo "❌ Error: Migration name is required"
    echo "Usage: ./create-migration.sh <migration_name>"
    echo "Example: ./create-migration.sh add_email_to_users"
    exit 1
fi

MIGRATION_NAME=$1
TIMESTAMP=$(date +%Y%m%d%H%M%S)
MIGRATION_NUMBER=$(ls -1 migrations/*.up.sql 2>/dev/null | wc -l | xargs printf "%06d")

UP_FILE="migrations/${MIGRATION_NUMBER}_${MIGRATION_NAME}.up.sql"
DOWN_FILE="migrations/${MIGRATION_NUMBER}_${MIGRATION_NAME}.down.sql"

# Create up migration file
cat > "$UP_FILE" << EOF
-- Migration: ${MIGRATION_NAME}
-- Created: $(date)

-- Add your migration SQL here

EOF

# Create down migration file
cat > "$DOWN_FILE" << EOF
-- Rollback: ${MIGRATION_NAME}
-- Created: $(date)

-- Add your rollback SQL here

EOF

echo "✅ Created migration files:"
echo "   - $UP_FILE"
echo "   - $DOWN_FILE"
