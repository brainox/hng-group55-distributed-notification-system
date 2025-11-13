# Email Service Integration Guide

## ✅ Integration Status

The Email Service is **READY** to work with the API Gateway. All message format mismatches have been fixed.

## 🔄 Message Flow

### 1. API Gateway → Email Queue
When you POST to `http://localhost:3000/api/v1/notifications/send`:

```json
{
  "notification_type": "email",
  "request_id": "200",
  "user_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "template_code": "welcome",
  "variables": {
    "name": "John",
    "link": "google.com" 
  },
  "priority": 1
}
```

### 2. RabbitMQ Queue Message Format
The API Gateway publishes this to `email.queue`:

```json
{
  "notification_id": "9bed08af-476f-469f-bdeb-0a2fb9c3ecc8",
  "notification_type": "email",
  "user_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "recipient": "user@example.com",
  "subject": "Test Notification",
  "body": "Hello John, this is a test notification!",
  "template_code": "welcome",
  "variables": {
    "name": "John",
    "link": "google.com"
  },
  "priority": 1,
  "metadata": {
    "timestamp": "2025-11-11T15:34:24.780Z",
    "retry_count": 0
  }
}
```

### 3. Email Service Processing
The Email Service:
1. ✅ Consumes message from `email.queue`
2. ✅ Checks idempotency (prevents duplicate sends)
3. ✅ Uses **pre-rendered** content from API Gateway (subject + body)
4. ✅ Falls back to Template Service if content not provided
5. ✅ Sends email via SMTP/SendGrid with circuit breaker
6. ✅ Retries on failure (exponential backoff)
7. ✅ Publishes status update back to RabbitMQ

## 🔧 Changes Made

### EmailMessage Model Updated
```go
type EmailMessage struct {
    NotificationID   string                 `json:"notification_id"`  // was: id
    NotificationType string                 `json:"notification_type"` // NEW
    UserID           string                 `json:"user_id"`          // NEW
    Recipient        string                 `json:"recipient"`
    Subject          string                 `json:"subject"`          // NEW - pre-rendered
    Body             string                 `json:"body"`             // NEW - pre-rendered
    TemplateCode     string                 `json:"template_code"`    // was: template_id
    Variables        map[string]interface{} `json:"variables"`
    Priority         int                    `json:"priority"`         // was: string
    Metadata         struct {
        Timestamp  string `json:"timestamp"`
        RetryCount int    `json:"retry_count"`
    } `json:"metadata"`                                               // NEW
}
```

### Consumer Logic Updated
- Uses `NotificationID` instead of `ID` and `CorrelationID`
- Uses `UserID` instead of `CorrelationID` for status updates
- Prefers pre-rendered content from API Gateway
- Falls back to Template Service if needed

## 🚀 How to Test

### Option 1: With Pre-rendered Content (Current API Gateway Behavior)
The API Gateway already renders templates, so your Email Service will:
1. Receive message with `subject` and `body` already populated
2. Skip Template Service call
3. Send email directly

```bash
# Start Email Service
cd services/email_service
go run cmd/server/main.go
```

### Option 2: With Template Service Integration
If you want to use the Template Service:

1. **Create a template** in Template Service:
```http
POST http://localhost:8081/api/v1/templates
Content-Type: application/json

{
  "template_key": "welcome",
  "name": "Welcome Email",
  "description": "Welcome email for new users",
  "template_type": "email",
  "subject": "Welcome {{name}}!",
  "body": "<h1>Hello {{name}}!</h1><p>Visit: <a href='{{link}}'>{{link}}</a></p>",
  "language": "en",
  "variables": ["name", "link"]
}
```

2. **Update API Gateway** to not render templates (remove the mock template rendering)

3. Email Service will automatically fetch from Template Service

## 📊 Status Updates

Email Service publishes status back to RabbitMQ:

```json
{
  "notification_id": "9bed08af-476f-469f-bdeb-0a2fb9c3ecc8",
  "user_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "status": "sent",
  "timestamp": "2025-11-11T15:35:00.000Z",
  "error": "",
  "provider": "smtp"
}
```

Status values:
- `sent` - Email sent successfully
- `failed` - Email failed to send

## 🔑 Environment Variables Required

```env
# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
EMAIL_QUEUE=email.queue
STATUS_QUEUE=email.status
EMAIL_WORKERS=3

# Template Service
TEMPLATE_SERVICE_URL=http://localhost:8081

# Redis
REDIS_URL=redis://localhost:6379

# Email Provider (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@yourcompany.com

# OR Email Provider (SendGrid)
EMAIL_PROVIDER=sendgrid
SENDGRID_API_KEY=your-sendgrid-api-key
SENDGRID_FROM_EMAIL=noreply@yourcompany.com
SENDGRID_FROM_NAME=Your Company

# Retry Configuration
MAX_RETRY_ATTEMPTS=5
RETRY_BASE_DELAY=1s

# Circuit Breaker
CIRCUIT_BREAKER_THRESHOLD=5
CIRCUIT_BREAKER_TIMEOUT=30s
```

## ✅ Checklist

- [x] Message format aligned with API Gateway
- [x] Field names updated (notification_id, template_code, etc.)
- [x] Pre-rendered content support added
- [x] Template Service fallback maintained
- [x] Status update format corrected
- [x] Idempotency checking works
- [x] Circuit breaker implemented
- [x] Retry logic with exponential backoff
- [x] Logging uses correct field names

## 🐛 Debugging

Watch logs while testing:
```bash
# Email Service logs will show:
# - "using pre-rendered content from API Gateway" (if using API Gateway rendering)
# - "rendered template from Template Service" (if fetching from Template Service)
# - "email sent successfully" (on success)
# - "failed to process email after retries" (on failure)
```

Check RabbitMQ Management UI:
- Queue: `email.queue` - should see messages consumed
- Queue: `email.status` - should see status updates published
- Exchange: `notifications.direct` - routing working

## 🎯 Next Steps

1. ✅ **Test with real SMTP credentials** - Update `.env` with actual email provider
2. ⏳ **Integrate User Service** - Replace mock recipient with real user lookup
3. ⏳ **Update API Gateway** - Remove mock template rendering, let Template Service handle it
4. ⏳ **Add monitoring** - Prometheus metrics for email success/failure rates
5. ⏳ **Add webhook callbacks** - Notify API Gateway of delivery status

## 📝 Notes

- API Gateway currently renders templates itself (mock rendering)
- Email Service is flexible: works with pre-rendered OR fetches from Template Service
- All field mappings are now correct
- Ready for production testing with real email provider
