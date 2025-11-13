# Template Service - Autonomous Deployment

## Architecture
This service is **fully autonomous** with its own embedded PostgreSQL database and Redis cache.

## Database
- **PostgreSQL 15**: Embedded within the service
- **Port**: 5433 (external), 5432 (internal)
- **Database Name**: `template_db`
- **Connection**: Only accessible within this service's network

## Running Independently

### Development (Standalone)
```bash
cd services/template_service
docker-compose up --build
```

This starts:
- Template Service API on port 8081
- Embedded PostgreSQL on port 5433
- Embedded Redis on port 6380

### Testing
```bash
# Health check
curl http://localhost:8081/health

# Create a template
curl -X POST http://localhost:8081/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "template_code": "test_template",
    "name": "Test Template",
    "subject": "Hello {{name}}",
    "body": "Welcome {{name}}!",
    "variables": ["name"]
  }'
```

## Running with Full System
When running via the root docker-compose.yml, this service's database is automatically included and networked with other services.

```bash
# From project root
docker-compose up template-service template-db
```

## Microservice Principles
✅ **Autonomous**: Can run independently without external database  
✅ **Loosely Coupled**: Own database, no shared state  
✅ **Self-Contained**: All dependencies bundled  
✅ **Isolated**: Separate network namespace  
