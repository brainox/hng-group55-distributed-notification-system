# Quick Start Guide

Get the complete notification system running in under 5 minutes.

## Prerequisites

- Docker & Docker Compose
- SMTP credentials (Gmail, SendGrid, or Mailtrap)

## Setup

1. **Configure SMTP**:
   ```bash
   cp .env.example .env
   # Edit .env with your SMTP credentials
   ```

2. **Start Services**:
   ```bash
   docker-compose up -d
   ```

3. **Verify Health**:
   ```bash
   # All services should be running
   docker-compose ps
   ```

4. **Test the System**:
   
   Open `test.http` in VS Code with the REST Client extension installed:
   
   ```bash
   # Install REST Client extension in VS Code
   # Search for "REST Client" by Huachao Mao
   ```
   
   Then open `test.http` and run the requests:
   
   - **Step 1**: Click "Send Request" above "### API Gateway Health" to verify all services
   - **Step 2**: Run "### Create Welcome Email Template" to create a template
   - **Step 3**: Run "### Send Email Notification with Template" to send a test email
   - **Step 4**: Update `your-email@example.com` with your real email address
   
   Alternative using curl:
   ```bash
   # Health check
   curl http://localhost:3000/api/v1/health
   
   # Create template
   curl -X POST http://localhost:8081/api/v1/templates \
     -H "Content-Type: application/json" \
     -d '{"template_key":"welcome","name":"Welcome Email","template_type":"email","subject":"Welcome {{name}}!","body":"Hello {{name}}, visit {{link}}","language":"en","is_active":true}'
   
   # Send notification
   curl -X POST http://localhost:3000/api/v1/notifications/send \
     -H "Content-Type: application/json" \
     -d '{"notification_type":"email","request_id":"test-001","user_id":"test-user","recipient":"your-email@example.com","template_code":"welcome","variables":{"name":"John","link":"https://example.com"},"priority":1}'
   ```

5. **Monitor Processing**:
   ```bash
   docker-compose logs -f email-service
   ```

## Services

| Service | Port | UI/API |
|---------|------|--------|
| API Gateway | 3000 | http://localhost:3000 |
| Template Service | 8081 | http://localhost:8081 |
| Email Service | 8082 | http://localhost:8082 |
| RabbitMQ | 5672, 15672 | http://localhost:15672 (guest/guest) |
| Redis | 6379 | - |
| PostgreSQL | 5432 | - |

## Architecture

```
User Request
    ↓
API Gateway (3000)
    ↓
RabbitMQ (email.queue)
    ↓
Email Service (8082)
    ↓
Template Service (8081) ← PostgreSQL (5432)
    ↓
SMTP Provider
    ↓
Email Delivered
```

## What's Next?

- **Full Testing Guide**: See [TESTING.md](./TESTING.md)
- **Integration Details**: See [services/email_service/INTEGRATION.md](./services/email_service/INTEGRATION.md)
- **Add More Templates**: Create password reset, verification, etc.
- **Scale Up**: Add more workers or containers
- **Monitor**: Check RabbitMQ UI at http://localhost:15672

## Troubleshooting

**Services not starting?**
```bash
docker-compose logs [service-name]
```

**Emails not sending?**
- Verify SMTP credentials in `.env`
- Gmail users: Use App Password (not regular password)
- Check Email Service logs: `docker-compose logs email-service`

**Messages stuck in queue?**
- Check RabbitMQ UI: http://localhost:15672
- Restart Email Service: `docker-compose restart email-service`

## Cleanup

```bash
# Stop services
docker-compose down

# Stop and remove all data
docker-compose down -v
```
