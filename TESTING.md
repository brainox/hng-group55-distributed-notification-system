# End-to-End Testing Guide

This guide walks through testing the complete notification system with all services running in Docker.

## Prerequisites

- Docker and Docker Compose installed
- Git repository cloned
- SMTP credentials (Gmail, SendGrid, Mailtrap, etc.)

## Setup

### 1. Configure Environment Variables

```bash
# Copy environment template
cp .env.example .env

# Edit .env with your SMTP credentials
nano .env
```

Update these values:
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@yourcompany.com
```

**Gmail Users**: Use App Password instead of regular password:
1. Enable 2FA on your Google account
2. Generate App Password at: https://myaccount.google.com/apppasswords
3. Use the generated password in `SMTP_PASSWORD`

### 2. Start All Services

```bash
# Build and start all services
docker-compose up -d

# Check all containers are running
docker-compose ps

# Expected output:
# notification-postgres          running (healthy)
# notification-rabbitmq          running (healthy)
# notification-redis             running (healthy)
# notification-api-gateway       running
# notification-template-service  running
# notification-email-service     running
```

### 3. View Logs

```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f api-gateway
docker-compose logs -f template-service
docker-compose logs -f email-service

# View last 100 lines
docker-compose logs --tail=100
```

## Testing Workflow

### Step 1: Verify Services Health

```bash
# API Gateway
curl http://localhost:3000/api/v1/health

# Template Service
curl http://localhost:8081/health

# Email Service
curl http://localhost:8082/health

# RabbitMQ Management UI
# Open browser: http://localhost:15672
# Login: guest / guest
```

### Step 2: Run Database Migrations

Template Service should auto-migrate on startup. Verify tables exist:

```bash
# Connect to PostgreSQL
docker exec -it notification-postgres psql -U postgres -d template_db

# List tables
\dt

# Expected tables:
# templates
# template_versions
# template_audit_logs

# Exit
\q
```

### Step 3: Create Email Template

```bash
# Create welcome template
curl -X POST http://localhost:8081/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "template_key": "welcome",
    "name": "Welcome Email",
    "template_type": "email",
    "subject": "Welcome {{name}} to Our Platform!",
    "body": "Hello {{name}},\n\nThank you for joining us. Get started by visiting: {{link}}\n\nBest regards,\nThe Team",
    "language": "en",
    "is_active": true,
    "metadata": {
      "category": "onboarding",
      "tags": ["welcome", "email"]
    }
  }'

# Verify template created
curl http://localhost:8081/api/v1/templates/welcome
```

### Step 4: Send Notification via API Gateway

```bash
# Send email notification
curl -X POST http://localhost:3000/api/v1/notifications/send \
  -H "Content-Type: application/json" \
  -d '{
    "notification_type": "email",
    "request_id": "test-001",
    "user_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "recipient": "recipient@example.com",
    "template_code": "welcome",
    "variables": {
      "name": "John Doe",
      "link": "https://yourapp.com/dashboard"
    },
    "priority": 1,
    "metadata": {
      "source": "manual_test"
    }
  }'

# Expected response:
# {
#   "notification_id": "uuid-here",
#   "status": "queued",
#   "message": "Notification queued successfully"
# }
```

### Step 5: Monitor Email Processing

```bash
# Watch Email Service logs
docker-compose logs -f email-service

# Expected log flow:
# 1. "Connected to RabbitMQ"
# 2. "Starting email workers: 3"
# 3. "Processing email message" (with notification_id)
# 4. "Fetching template: welcome" (if not pre-rendered)
# 5. "Template fetched successfully"
# 6. "Rendering email with variables"
# 7. "Sending email to: recipient@example.com"
# 8. "Email sent successfully"
# 9. "Publishing status update: sent"
```

### Step 6: Check RabbitMQ

Open RabbitMQ Management UI: http://localhost:15672 (guest/guest)

1. **Queues Tab**:
   - `email.queue`: Should show messages consumed
   - `push.queue`: Empty (not used in this test)
   - `failed.queue`: Should be empty (no failures)

2. **Connections Tab**: Should show Email Service consumer connected

3. **Exchanges Tab**: `notifications.direct` should show message routes

### Step 7: Verify Email Delivery

- **Check your SMTP provider's dashboard** (Gmail, SendGrid, etc.)
- **Check recipient inbox** (if using real email)
- **Mailtrap Users**: Check https://mailtrap.io/inboxes

### Step 8: Test Pre-rendered Content

API Gateway can send pre-rendered emails (bypasses Template Service):

```bash
curl -X POST http://localhost:3000/api/v1/notifications/send \
  -H "Content-Type: application/json" \
  -d '{
    "notification_type": "email",
    "request_id": "test-002",
    "user_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "recipient": "recipient@example.com",
    "subject": "Custom Subject",
    "body": "This is a pre-rendered email body. No template needed!",
    "priority": 1
  }'
```

Email Service will use the provided subject/body directly.

## Testing Scenarios

### Scenario 1: Template-based Email (Happy Path)

1. Create template in Template Service
2. Send notification with `template_code` and `variables`
3. Email Service fetches template and renders
4. Email sent successfully

### Scenario 2: Pre-rendered Email (API Gateway renders)

1. API Gateway renders template internally
2. Send notification with `subject` and `body` already populated
3. Email Service uses provided content directly
4. Email sent successfully

### Scenario 3: Template Not Found

1. Send notification with non-existent `template_code`
2. Email Service attempts to fetch template
3. Template Service returns 404
4. Email Service logs error
5. Message moves to failed queue

### Scenario 4: SMTP Failure

1. Stop SMTP server or use invalid credentials
2. Send notification
3. Email Service attempts to send
4. SMTP fails, triggers retry logic
5. After max retries, message moves to failed queue

### Scenario 5: High Load

```bash
# Send 100 notifications
for i in {1..100}; do
  curl -X POST http://localhost:3000/api/v1/notifications/send \
    -H "Content-Type: application/json" \
    -d "{
      \"notification_type\": \"email\",
      \"request_id\": \"load-test-$i\",
      \"user_id\": \"f47ac10b-58cc-4372-a567-0e02b2c3d479\",
      \"recipient\": \"recipient@example.com\",
      \"template_code\": \"welcome\",
      \"variables\": {\"name\": \"User $i\", \"link\": \"https://example.com\"},
      \"priority\": 1
    }"
done

# Monitor Email Service processing
docker-compose logs -f email-service | grep "Email sent successfully"
```

## Troubleshooting

### Services Not Starting

```bash
# Check logs for specific service
docker-compose logs template-service
docker-compose logs email-service

# Common issues:
# - Port conflicts: Change ports in docker-compose.yml
# - Build errors: Check Dockerfile and dependencies
# - Health check failures: Increase timeout/retries
```

### Template Service Can't Connect to PostgreSQL

```bash
# Verify PostgreSQL is running
docker-compose ps postgres

# Check PostgreSQL logs
docker-compose logs postgres

# Test connection
docker exec notification-postgres psql -U postgres -d template_db -c "SELECT 1;"
```

### Email Service Not Consuming Messages

```bash
# Check RabbitMQ connection
docker-compose logs email-service | grep "RabbitMQ"

# Verify queue exists
curl -u guest:guest http://localhost:15672/api/queues/%2F/email.queue

# Restart Email Service
docker-compose restart email-service
```

### Emails Not Sending

```bash
# Check SMTP configuration
docker-compose logs email-service | grep "SMTP"

# Common SMTP issues:
# - Gmail: Must use App Password (not regular password)
# - Port blocked: Try port 465 (SSL) or 2525 (alternative)
# - Authentication failed: Verify credentials
# - Firewall: Check if SMTP ports are open

# Test SMTP connectivity from container
docker exec notification-email-service nc -zv smtp.gmail.com 587
```

### Messages Going to Failed Queue

```bash
# Check failed queue in RabbitMQ UI
# http://localhost:15672/#/queues/%2F/failed.queue

# Or via API
curl -u guest:guest http://localhost:15672/api/queues/%2F/failed.queue

# Common reasons:
# - Template not found
# - Invalid email format
# - SMTP failures after max retries
# - Invalid message format
```

### High Memory/CPU Usage

```bash
# Check resource usage
docker stats

# Reduce Email Service workers
# Edit docker-compose.yml: EMAIL_WORKERS=1

# Restart services
docker-compose restart email-service
```

## Performance Monitoring

### RabbitMQ Metrics

- **Message Rate**: http://localhost:15672/#/queues (Messages/sec)
- **Consumer Utilization**: Should be >50% under load
- **Queue Length**: Should stay near 0 under normal load

### Redis Metrics

```bash
# Connect to Redis CLI
docker exec -it notification-redis redis-cli

# Check memory usage
INFO memory

# Check keys count
DBSIZE

# Monitor commands in real-time
MONITOR
```

### Email Service Metrics

```bash
# Count successful sends
docker-compose logs email-service | grep "Email sent successfully" | wc -l

# Count failures
docker-compose logs email-service | grep "Failed to send email" | wc -l

# Average processing time (check logs for timing)
docker-compose logs email-service | grep "Processing completed in"
```

## Cleanup

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (deletes all data)
docker-compose down -v

# Remove images
docker-compose down --rmi all
```

## Integration Status

✅ **API Gateway → RabbitMQ**: Messages published successfully  
✅ **RabbitMQ → Email Service**: Messages consumed successfully  
✅ **Email Service → Template Service**: Templates fetched successfully  
✅ **Email Service → SMTP**: Emails sent successfully  
✅ **Status Updates → Redis**: Status tracking working  

## Next Steps

1. **Add More Templates**: Create templates for different use cases
2. **Monitor Production**: Set up logging/monitoring (ELK, Datadog, etc.)
3. **Scale Services**: Add more Email Service workers or containers
4. **Add Push Service**: Implement push notification service
5. **Add SMS Service**: Implement SMS notification service
6. **CI/CD**: Set up automated testing and deployment
7. **Security**: Add authentication, rate limiting, encryption

## Support

For issues or questions:
1. Check service logs: `docker-compose logs [service-name]`
2. Review INTEGRATION.md in email_service
3. Check RabbitMQ UI: http://localhost:15672
4. Verify environment variables in .env file
