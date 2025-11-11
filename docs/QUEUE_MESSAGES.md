# Queue Message Formats

## Email Queue Message

**Queue**: `email.queue`  
**Routing Key**: `email`
```json
{
  "notification_id": "notif-abc-123",
  "notification_type": "email",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "recipient": "user@example.com",
  "subject": "Welcome!",
  "body": "Hello John, welcome to our app!",
  "variables": {
    "name": "John",
    "link": "https://example.com"
  },
  "request_id": "req-12345-abc",
  "template_code": "welcome-email",
  "priority": 1,
  "metadata": {
    "timestamp": "2025-11-09T10:00:00Z",
    "retry_count": 0
  }
}
```

## Push Queue Message

**Queue**: `push.queue`  
**Routing Key**: `push`
```json
{
  "notification_id": "notif-def-456",
  "notification_type": "push",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "recipient": "fcm-token-abc123",
  "title": "New Message",
  "body": "You have 3 unread messages",
  "template_code": "new_message_push",
  "variables": {
    "name": "John Doe",
    "link": "https://example.com/messages"
  },
  "request_id": "req-67890-def",
  "priority": 3,
  "metadata": {
    "timestamp": "2025-11-09T10:00:00Z",
    "retry_count": 0
  }
}
```

## Connection Details

- **AMQP URL**: `amqp://localhost:5672`
- **Management UI**: http://localhost:15672
- **Credentials**: guest / guest
```