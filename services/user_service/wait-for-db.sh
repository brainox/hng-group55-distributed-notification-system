#!/bin/sh
# Wait until Postgres is ready
echo "Waiting for database..."
until nc -z -v -w30 $DATABASE_URL 
do
  echo "Postgres is unavailable - sleeping"
  sleep 2
done
echo "Postgres is up - running migrations"
alembic upgrade head
echo "Starting FastAPI"
uvicorn app.main:app --host 0.0.0.0 --port 8000
