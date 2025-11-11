# Push Service (Python)

This is a minimal push service implementation in Python. It supports:
- Reading messages from a Redis-backed queue
- Sending push notifications via:
  - Firebase Cloud Messaging (FCM) - legacy HTTP endpoint (or use firebase_admin)
  - OneSignal REST API
  - Web Push (VAPID) via pywebpush
- Token/subscription validation helpers
- Rich notifications (title, body, image, and optional click URL / data)

## Files
- `push_service/queue.py` - small Redis queue wrapper
- `push_service/validator.py` - token and subscription validators
- `push_service/senders.py` - FCM / OneSignal / WebPush sender implementations
- `push_service/worker.py` - example worker that consumes the queue and dispatches messages
- `requirements.txt` - required packages

## Example message (push_queue payload)
```json
{
  "provider": "fcm",
  "target": "FCM_DEVICE_TOKEN",
  "title": "Hello!",
  "body": "You have a new message.",
  "image": "https://example.com/image.png",
  "url": "https://example.com/app",
  "data": {"foo":"bar"}
}
```

## How to run
1. Install dependencies: `pip install -r requirements.txt`
2. Configure senders (FCM server key, OneSignal app/rest key, VAPID keys)
3. Push JSON messages into Redis list `push_queue` or change queue name
4. Start worker:
```py
from push_service.queue import RedisQueue
from push_service.senders import FCMSender, OneSignalSender, WebPushSender
from push_service.worker import PushWorker

q = RedisQueue(redis_url='redis://localhost:6379/0', queue_name='push_queue')
fcm = FCMSender(server_key='YOUR_FCM_SERVER_KEY')
one = OneSignalSender(app_id='YOUR_APP_ID', rest_api_key='YOUR_REST_KEY')
web = WebPushSender(vapid_private_key='YOUR_PRIV_KEY', vapid_claims={'sub':'mailto:you@example.com'})

worker = PushWorker(q, fcm_sender=fcm, onesignal_sender=one, webpush_sender=web)
worker.run_forever()
```

## Notes
- This is a starting point. In production, consider:
  - Using FCM HTTP v1 with OAuth service accounts or firebase_admin SDK
  - Retries, dead-letter queue, rate-limiting, metrics, and structured logging
  - Securely storing credentials and rotating keys
  - Using a robust queue system (RabbitMQ, AWS SQS, Google PubSub)
