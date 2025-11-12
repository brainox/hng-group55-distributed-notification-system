#!/bin/bash
set -e

# Create users_db database
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    SELECT 'CREATE DATABASE users_db'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'users_db')\gexec
EOSQL

echo "Database users_db created successfully!"
