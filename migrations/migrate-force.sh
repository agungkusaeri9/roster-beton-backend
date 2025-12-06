#!/bin/bash

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

# Change to project root directory
cd "$PROJECT_ROOT"

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Set default values if not set
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-go_arch}

# PostgreSQL connection string
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Check if version argument is provided
if [ -z "$1" ]; then
    echo "❌ Error: Version number is required"
    echo "Usage: ./migrate-force.sh <version>"
    echo "Example: ./migrate-force.sh 1"
    exit 1
fi

VERSION=$1

# Find migrate binary
MIGRATE_CMD="migrate"
if ! command -v migrate &> /dev/null; then
    if [ -f "$HOME/go/bin/migrate" ]; then
        MIGRATE_CMD="$HOME/go/bin/migrate"
    elif [ -f "/usr/local/go/bin/migrate" ]; then
        MIGRATE_CMD="/usr/local/go/bin/migrate"
    else
        echo "❌ Error: migrate command not found"
        exit 1
    fi
fi

# Force migration to specific version
echo "🔧 Forcing migration to version $VERSION..."
$MIGRATE_CMD -path ./migrations -database "$DATABASE_URL" force $VERSION

