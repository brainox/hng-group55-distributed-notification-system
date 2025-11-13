#!/bin/sh
echo "Running Alembic migrations..."
alembic upgrade head || { echo "Migration failed"; exit 1; }
echo "Starting FastAPI..."
uvicorn app.main:app --host 0.0.0.0 --port 8000
