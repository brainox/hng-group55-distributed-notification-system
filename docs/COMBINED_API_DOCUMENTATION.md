Distributed Notification System — Combined API Documentation

This document contains the full API documentation for the entire distributed notification system consisting of:

API Gateway
User Service
Template Service
Email Service
Push Service

1.API Gateway
📌 Base URL
/api/v1

----api/v1/notifications/

Send an email request.

Request

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

2.  User Service
    📌 Base URL
    `----v1/users/

    POST /users/Login
    ----- v1/users/login

    ```json
    {
      "email": "zumma34@gmail.com",
      "password": "zumma126"
    }
    ```

    POST /users/Register
    Registers Users
    ------v1/users/register

    ```json
    {
      "name": "zumma jekob",
      "email": "zumma34@gmail.com",
      "password": "zumma126",
      "push_token": "token_123",
      "preferences": {
        "email": true,
        "push": true
      }
    }
    ```

    GET /users/{users_id}
    ----/v1/users/{users_id}
    Get users profile

    PUT /users/{user_id}/preferences
    -----/v1/users/{users_id}/preferences
    set notification preference

    ```json
    {
      "email": false,
      "push": true
    }
    ```

3.  📝 Template Service
    📌 Base URL
    -------/v1/templates

    POST /templates
    creates a template

    ```json
    {
      "template_key": "template_key_1",
      "name": "Welcome to HNG",
      "description": "Welcome to HNG template message description",
      "template_type": "email",
      "subject": "Look who got into HNG! {{name}}!",
      "body": "Hello {{name}}, welcome to HNG! HURRAY!! Visit: {{link}}",
      "language": "en",
      "variables": ["name", "link"]
    }
    ```

    GET /templates
    Fetches all templates

4.  📧 Email Service
    This service listens on:
    email.queue

        Responsibilities:
        Consume email messages
        Replace template variables
        Send email via SMTP/SendGrid
        Update status back to gateway

    No external APIs (internal only)

5.  🔔 Push Service
    GET /notifications
    Retrieves a comprehensive list of all stored push notification messages, including their current statuses and metadata.

    GET /notifications/{notification_id}
    Retrieves the details of a single push notification message using its unique identifier.

6.  📨 Queue Messages
    Sample queue message (email):

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


7. 📦 Unified Response Format
      {
      "success": true,
      "data": {},
      "error": null,
      "message": "string",
      "meta": {
      "total": 0,
      "limit": 10,
      "page": 1,
      "total_pages": 1,
      "has_next": false,
      "has_previous": false
      }
      }

8. 🧩 Data Models (snake_case)
      User
      id: uuid
      email: string
      password_hash: string
      push_token: string
      preferences: JSON

      Template
      id: uuid
      code: string
      language: string
      content: text
      version: integer
      created_at: timestamp

9. 💚 Health Checks

Every service MUST expose:

GET /health

Response:

{ "status": "OK" }

10. 🏗️ System Architecture Overview

API Gateway → RabbitMQ → (Email/Push Services)
User & Template Services via REST
Redis for caching preferences & tokens
PostgreSQL for persistent storage
Dead-letter queue for failed messages
```
