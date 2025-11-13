# Autonomous Microservices Architecture

## Overview
This notification system now follows **true microservice principles** where each service is fully autonomous with its own embedded database.

## Architecture Changes

### Before (Shared Database - Anti-Pattern)
```
Root docker-compose.yml
├── postgres (shared)
│   ├── template_db
│   └── users_db
├── template-service → postgres
└── user-service → postgres
```
**Problems:**
- ❌ Tight coupling between services
- ❌ Single point of failure
- ❌ Cannot deploy services independently
- ❌ Database schema conflicts
- ❌ Difficult to scale individual services

### After (Autonomous Services - Best Practice)
```
services/
├── template_service/
│   ├── docker-compose.yml
│   ├── postgres (embedded)
│   ├── redis (embedded)
│   └── service code
│
└── user_service/
    ├── docker-compose.yaml
    ├── postgres (embedded)
    └── service code
```
**Benefits:**
- ✅ Complete service autonomy
- ✅ Independent deployment
- ✅ Loose coupling
- ✅ Service-specific optimizations
- ✅ Fault isolation
- ✅ Individual scaling

## Service Independence

### Template Service
**Location**: `services/template_service/`

**Runs Independently**:
```bash
cd services/template_service
docker-compose up
```

**Includes**:
- Template Service API (port 8081)
- PostgreSQL 15 (port 5433, database: template_db)
- Redis 7 (port 6380)
- Own network: `template-network`

**Can be deployed separately** to its own infrastructure.

### User Service
**Location**: `services/user_service/`

**Runs Independently**:
```bash
cd services/user_service
docker-compose up
```

**Includes**:
- User Service API (port 3001)
- PostgreSQL 15 (port 5434, database: users_db)
- Own network: `user-network`
- Alembic migrations

**Can be deployed separately** to its own infrastructure.

## Full System Orchestration

The root `docker-compose.yml` now **orchestrates** all services without owning databases:

```yaml
# Root docker-compose.yml uses "extends" pattern
services:
  template-service:
    extends:
      file: ./services/template_service/docker-compose.yml
      service: template-service
  
  template-db:
    extends:
      file: ./services/template_service/docker-compose.yml
      service: template-db
  
  user-service:
    extends:
      file: ./services/user_service/docker-compose.yaml
      service: user_service
  
  user-db:
    extends:
      file: ./services/user_service/docker-compose.yaml
      service: db
```

This allows:
- Services to be developed independently
- Each service owns its database lifecycle
- Root compose is for **integration testing** only
- Production: deploy each service separately

## Running Options

### Option 1: Independent Development
Work on one service without affecting others:
```bash
# Work on template service only
cd services/template_service
docker-compose up

# Work on user service only
cd services/user_service
docker-compose up
```

### Option 2: Full System Integration
Test all services together:
```bash
# From project root
docker-compose up
```

### Option 3: Selective Services
Run only what you need:
```bash
# Just user service + API Gateway
docker-compose up user-service user-db api-gateway

# Just template service + email service
docker-compose up template-service template-db email-service
```

## Port Mapping

### Shared Infrastructure
- RabbitMQ: 5672, 15672 (management)
- Redis (shared): 6379

### Template Service
- API: 8081
- PostgreSQL: 5433
- Redis: 6380

### User Service
- API: 3001 (internal: 8000)
- PostgreSQL: 5434

### Other Services
- API Gateway: 3000
- Email Service: 8082

## Database Access

Each service's database is **isolated by default**:

```bash
# Access template database
docker exec -it template-db psql -U postgres -d template_db

# Access user database  
docker exec -it user-db psql -U postgres -d users_db
```

## Microservice Principles Achieved

| Principle | Implementation |
|-----------|----------------|
| **Autonomous** | Each service has own DB, can deploy independently |
| **Loosely Coupled** | No shared database, services communicate via APIs |
| **Single Responsibility** | Each service manages only its domain data |
| **Fault Isolation** | Template DB down ≠ User Service down |
| **Independent Scaling** | Scale services based on their own load |
| **Technology Freedom** | Services can use different DB versions/configs |
| **Data Sovereignty** | Complete control over schema and migrations |

## Production Deployment

In production, deploy each service to separate infrastructure:

```bash
# Deploy template service to its own server/cluster
cd services/template_service
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Deploy user service to its own server/cluster
cd services/user_service
docker-compose -f docker-compose.yaml -f docker-compose.prod.yaml up -d
```

Each service can be:
- Scaled independently
- Updated without affecting others
- Monitored separately
- Backed up on its own schedule
- Optimized for its specific workload

## Migration Guide

If you previously had data in the shared database:

```bash
# 1. Export data from old setup
docker exec old-postgres pg_dump -U postgres template_db > template_backup.sql
docker exec old-postgres pg_dump -U postgres users_db > users_backup.sql

# 2. Start new services
docker-compose up -d

# 3. Import data
docker exec -i template-db psql -U postgres -d template_db < template_backup.sql
docker exec -i user-db psql -U postgres -d users_db < users_backup.sql
```

## Summary

This architecture transformation moves from a **monolithic database approach** to **true microservices**, where:

- ✅ Each service is a complete, deployable unit
- ✅ Services can be developed, tested, and deployed independently
- ✅ Failures are isolated
- ✅ Scaling is service-specific
- ✅ Teams can work autonomously on their services

This is the **industry best practice** for distributed systems and microservice architectures.
