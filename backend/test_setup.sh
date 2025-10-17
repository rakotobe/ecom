#!/bin/bash
# Test database setup script

set -e

DB_NAME="ecom_test"
DB_USER="postgres"
DB_PASSWORD="postgres"
DB_HOST="localhost"
DB_PORT="5432"

echo "Setting up test database..."

# Check if PostgreSQL is running
if ! pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER > /dev/null 2>&1; then
    echo "PostgreSQL is not running on $DB_HOST:$DB_PORT"
    echo "Please start PostgreSQL or use Docker:"
    echo "  docker run --name postgres-test -e POSTGRES_PASSWORD=$DB_PASSWORD -p $DB_PORT:5432 -d postgres:16-alpine"
    exit 1
fi

# Drop and recreate test database
echo "Dropping existing test database (if any)..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "DROP DATABASE IF EXISTS $DB_NAME;" 2>/dev/null || true

echo "Creating test database..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $DB_NAME;"

echo "Test database setup complete!"
echo "Run integration tests with: go test -v ./infrastructure/... ./api/..."
