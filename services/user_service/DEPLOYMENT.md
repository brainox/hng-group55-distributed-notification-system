# User Service - Autonomous Deployment

## Architecture
This service is **fully autonomous** with its own embedded PostgreSQL database.

## Database
- **PostgreSQL 15**: Embedded within the service
- **Port**: 5434 (external), 5432 (internal)
- **Database Name**: `users_db`
- **Connection**: Only accessible within this service's network

## Running Independently

### Development (Standalone)
```bash
cd services/user_service
docker-compose up --build
```

This starts:
- User Service API on port 3001 (mapped from internal 8000)
- Embedded PostgreSQL on port 5434
- Automatic database migrations via Alembic

### Testing
```bash
# Health check
curl http://localhost:3001/v1/health

# Register a user
curl -X POST http://localhost:3001/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123",
    "preferences": {"email": true, "push": true, "sms": false}
  }'

# Get user by ID
curl http://localhost:3001/v1/users/{user_id}
```

## Running with Full System
When running via the root docker-compose.yml, this service's database is automatically included and networked with other services.

```bash
# From project root
docker-compose up user-service user-db
```

## Database Migrations
Migrations are managed by Alembic and run automatically on startup via `start.sh`.

```bash
# Inside container
alembic upgrade head

# Create new migration
alembic revision --autogenerate -m "description"
```

## Microservice Principles
✅ **Autonomous**: Can run independently without external database  
✅ **Loosely Coupled**: Own database, no shared state  
✅ **Self-Contained**: All dependencies bundled  
✅ **Isolated**: Separate network namespace  
✅ **Data Sovereignty**: Complete control over user data schema and migrations  
