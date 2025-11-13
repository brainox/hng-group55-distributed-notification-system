# Email Service

The Email Service is a microservice responsible for consuming email notification messages from RabbitMQ, fetching templates from the Template Service, rendering them with dynamic variables, and sending emails via SMTP or SendGrid.

## Features

- **RabbitMQ Consumer**: Processes messages from `email.queue` with a worker pool (10 concurrent workers)
- **Template Integration**: Fetches templates from Template Service with Redis caching (10-minute TTL)
- **Variable Rendering**: Supports `{{variable}}` syntax for dynamic content
- **Multi-Provider Support**: SMTP (Gmail) and SendGrid email providers
- **Idempotency**: Prevents duplicate email sends using Redis (24-hour TTL)
- **Circuit Breaker**: Protects against cascading failures (opens after 5 consecutive failures)
- **Retry Logic**: Exponential backoff (1s, 2s, 4s, 8s, 16s) with intelligent permanent error detection
- **Status Updates**: Publishes success/failure status to `notification.status.queue`
- **Health Checks**: HTTP endpoint for monitoring service and dependencies
- **Graceful Shutdown**: Ensures in-flight messages are processed before shutdown

## Architecture

```
┌─────────────────┐
│   RabbitMQ      │
│  email.queue    │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────┐
│   Email Service (Port 8082) │
│  ┌──────────────────────┐   │
│  │  Worker Pool (10)    │   │
│  └──────────────────────┘   │
│           │                  │
│           ▼                  │
│  ┌──────────────────────┐   │
│  │  Idempotency Check   │   │
│  └──────────────────────┘   │
│           │                  │
│           ▼                  │
│  ┌──────────────────────┐   │
│  │  Template Fetcher    │───┼──► Template Service
│  │  (with Cache)        │   │
│  └──────────────────────┘   │
│           │                  │
│           ▼                  │
│  ┌──────────────────────┐   │
│  │  Template Renderer   │   │
│  └──────────────────────┘   │
│           │                  │
│           ▼                  │
│  ┌──────────────────────┐   │
│  │  Email Sender        │───┼──► SMTP/SendGrid
│  │  (Circuit Breaker)   │   │
│  └──────────────────────┘   │
│           │                  │
│           ▼                  │
│  ┌──────────────────────┐   │
│  │  Status Publisher    │   │
│  └──────────────────────┘   │
└──────────────┬──────────────┘
               │
               ▼
       ┌─────────────────┐
       │   RabbitMQ      │
       │ status.queue    │
       └─────────────────┘
```

## Message Format

### Input: Email Queue Message
```json
{
  "id": "unique-message-id",
  "correlation_id": "request-correlation-id",
  "recipient": "user@example.com",
  "template_id": "welcome_email",
  "variables": {
    "user_name": "John Doe",
    "activation_link": "https://example.com/activate/token"
  },
  "priority": "high",
  "scheduled_at": "2025-01-20T15:00:00Z"
}
```

### Output: Status Queue Message
```json
{
  "notification_id": "unique-message-id",
  "correlation_id": "request-correlation-id",
  "status": "sent",
  "provider": "smtp",
  "timestamp": "2025-01-20T15:00:05Z",
  "error": ""
}
```

Status values: `sent`, `failed`

## Configuration

Create a `.env` file in the service root:

```env
# Server
SERVER_PORT=8082
LOG_LEVEL=info

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
EMAIL_QUEUE_NAME=email.queue
STATUS_QUEUE_NAME=notification.status.queue
WORKER_COUNT=10

# Redis
REDIS_URL=localhost:6379
IDEMPOTENCY_TTL=86400

# Template Service
TEMPLATE_SERVICE_URL=http://localhost:8081/api/v1

# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@example.com

# SendGrid (Optional)
SENDGRID_API_KEY=

# Retry Configuration
MAX_RETRIES=5
RETRY_BACKOFF=1s

# Circuit Breaker Configuration
CIRCUIT_BREAKER_THRESHOLD=5
CIRCUIT_BREAKER_TIMEOUT=30s
```

## Getting Started

### Prerequisites

- Go 1.21+
- RabbitMQ running on port 5672
- Redis running on port 6379
- Template Service running on port 8081

### Installation

```bash
cd services/email_service

# Install dependencies
go mod download

# Run the service
go run cmd/server/main.go
```

### Docker

```bash
# Build image
docker build -t email-service:latest .

# Run container
docker run -d \
  --name email-service \
  --env-file .env \
  -p 8082:8082 \
  email-service:latest
```

## API Endpoints

### Health Check
```http
GET /health
```

**Response (200 OK)**
```json
{
  "status": "healthy",
  "checks": {
    "rabbitmq": "healthy",
    "redis": "healthy",
    "template_service": "healthy"
  },
  "timestamp": "2025-01-20T15:00:00Z"
}
```

**Response (503 Service Unavailable)**
```json
{
  "status": "unhealthy",
  "checks": {
    "rabbitmq": "healthy",
    "redis": "unhealthy: connection refused",
    "template_service": "healthy"
  },
  "timestamp": "2025-01-20T15:00:00Z"
}
```

## Email Provider Setup

### Gmail SMTP

1. Enable 2-factor authentication on your Google account
2. Generate an App Password: https://myaccount.google.com/apppasswords
3. Use the app password in `SMTP_PASSWORD`

### SendGrid

1. Create a SendGrid account
2. Generate an API key with "Mail Send" permission
3. Set `SENDGRID_API_KEY` in your environment
4. The service will automatically use SendGrid if the API key is present

## Retry Logic

The service implements intelligent retry logic:

**Retryable Errors:**
- Network errors
- Timeout errors
- Rate limit errors (429)
- Server errors (5xx)

**Non-Retryable Errors (Permanent):**
- Invalid email address
- Template not found
- Authentication failures (401, 403)
- Bad request (400)

**Backoff Schedule:**
- Attempt 1: Immediate
- Attempt 2: 1 second
- Attempt 3: 2 seconds
- Attempt 4: 4 seconds
- Attempt 5: 8 seconds
- Attempt 6: 16 seconds (max)

## Circuit Breaker

Protects the service from cascading failures:

- **Threshold**: Opens after 5 consecutive failures
- **Timeout**: Remains open for 30 seconds
- **Half-Open**: Allows 1 request to test recovery
- **Closed**: Normal operation

## Idempotency

Prevents duplicate email sends:

- Uses Redis with key format: `email:processed:{message_id}`
- TTL: 24 hours
- If a message ID is already processed, it's acknowledged without resending

## Monitoring

### Logs

The service uses structured logging (Zap):

```json
{
  "level": "info",
  "timestamp": "2025-01-20T15:00:00Z",
  "message": "processing email",
  "id": "msg-123",
  "correlation_id": "req-456",
  "recipient": "user@example.com",
  "template_id": "welcome_email"
}
```

### Metrics

Key events logged:
- Message consumption
- Idempotency checks
- Template fetching
- Email sending
- Retry attempts
- Circuit breaker state changes
- Status publishing

## Troubleshooting

### Message Not Processing

1. Check RabbitMQ connection:
   ```bash
   curl http://localhost:8082/health
   ```

2. Verify message format in queue

3. Check logs for errors:
   ```bash
   docker logs email-service
   ```

### Template Not Found

1. Verify Template Service is running:
   ```bash
   curl http://localhost:8081/health
   ```

2. Check template exists:
   ```bash
   curl http://localhost:8081/api/v1/templates/key/welcome_email
   ```

### SMTP Authentication Failed

1. Verify Gmail App Password is correct
2. Check 2FA is enabled
3. Test SMTP connection:
   ```bash
   openssl s_client -starttls smtp -connect smtp.gmail.com:587
   ```

### Circuit Breaker Open

1. Check email provider status
2. Verify credentials
3. Wait 30 seconds for automatic recovery
4. Review logs for failure patterns

### High Memory Usage

1. Reduce `WORKER_COUNT` (default: 10)
2. Decrease `IDEMPOTENCY_TTL`
3. Monitor Redis memory usage

## Development

### Project Structure

```
email_service/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration loading
│   ├── models/
│   │   └── email.go             # Data structures
│   ├── queue/
│   │   ├── consumer.go          # RabbitMQ consumer
│   │   └── publisher.go         # Status publisher
│   ├── sender/
│   │   ├── interface.go         # Email sender interface
│   │   ├── smtp.go              # SMTP implementation
│   │   └── sendgrid.go          # SendGrid implementation
│   ├── template/
│   │   ├── client.go            # Template Service HTTP client
│   │   └── renderer.go          # Variable substitution
│   ├── idempotency/
│   │   └── checker.go           # Duplicate detection
│   ├── circuit/
│   │   └── breaker.go           # Circuit breaker wrapper
│   ├── retry/
│   │   └── handler.go           # Retry logic
│   └── health/
│       └── checker.go           # Health checks
├── pkg/
│   └── logger/
│       └── logger.go            # Structured logging
├── .env.example                 # Environment template
├── go.mod                       # Go dependencies
└── README.md                    # This file
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/retry
```

### Adding a New Email Provider

1. Implement the `EmailSender` interface:
   ```go
   type EmailSender interface {
       Send(to, subject, body string) error
       GetProviderName() string
   }
   ```

2. Create a new file in `internal/sender/`

3. Update `cmd/server/main.go` to use your provider

## License

MIT
