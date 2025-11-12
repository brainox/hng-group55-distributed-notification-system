#!/usr/bin/env python3
import re
import time
import json
from typing import Dict, Any

import pika
from sqlalchemy.orm import Session
from pika.exceptions import AMQPConnectionError, ChannelClosedByBroker

from db import models
from logger import Logger
from db.database import SessionLocal
from service.senders import FCMSender
from db.crud import increment_retry_count


RABBITMQ_HOST = "localhost"
MAIN_QUEUE = "push_queue"
RETRY_QUEUE = "push_queue_retry"
DLQ_QUEUE = "push_queue_dlq"

MAX_RETRIES = 3
RETRY_DELAY_MS = 5000  # 5 seconds before retry


sender = FCMSender(
    credential_path="distributed-systems-3349d-firebase-adminsdk-fbsvc-b84f851127.json"
)


def handle_message(msg: Dict[str, Any]) -> bool:
    title = msg.get("title", "")
    body = msg.get("body", "")
    image = msg.get("image")
    url = msg.get("url")
    data = msg.get("data")
    token = msg.get("token")

    def validate_fcm_token(token: str) -> bool:
        return bool(token) and bool(re.fullmatch(r"[A-Za-z0-9:_\\-]{20,400}", token))

    if not validate_fcm_token(str(token)):
        Logger.warning("Invalid FCM token: %s", token)
        return False
    if not sender:
        Logger.error("FCM sender not configured")
        return False

    resp = sender.send(
        str(token), title, body, image=image, data=data, click_action=url
    )
    Logger.info("FCM response: %s", resp)
    return True if resp else False


def setup_queues(channel):
    """Declare main, retry, and DLQ queues with proper bindings."""
    # DLQ (no further retries)
    channel.queue_declare(queue=DLQ_QUEUE, durable=True)

    # Retry queue with TTL that sends messages back to main after delay
    channel.queue_declare(
        queue=RETRY_QUEUE,
        durable=True,
        arguments={
            "x-dead-letter-exchange": "",
            "x-dead-letter-routing-key": MAIN_QUEUE,
            "x-message-ttl": RETRY_DELAY_MS,
        },
    )

    # Main queue — failed messages go to retry queue
    channel.queue_declare(
        queue=MAIN_QUEUE,
        durable=True,
        arguments={
            "x-dead-letter-exchange": "",
            "x-dead-letter-routing-key": RETRY_QUEUE,
        },
    )


def process_message(ch, method, properties, body):
    db = SessionLocal()
    try:
        data = json.loads(body.decode("utf-8"))
        Logger.info(f"Processing message {data.get('id')} - {data.get('title')}")

        success = handle_message(data)
        msg = (
            db.query(models.PushMessage)
            .filter(models.PushMessage.id == data["id"])
            .first()
        )

        if msg:
            msg.status = (
                models.StatusEnum.success if success else models.StatusEnum.failed
            )
            db.commit()
            Logger.info(f"Message {msg.id} updated -> {msg.status}")
        else:
            Logger.warning(f"Message {data['id']} not found in DB")

        if not success:
            retry_count = msg.retry_count if msg and msg.retry_count else 0

            if retry_count >= MAX_RETRIES:
                Logger.error(f"Message {data['id']} failed permanently -> DLQ")
                ch.basic_publish(
                    exchange="",
                    routing_key=DLQ_QUEUE,
                    body=body,
                    properties=pika.BasicProperties(
                        delivery_mode=2,
                        headers={"x-original-error": "max retries exceeded"},
                    ),
                )
                ch.basic_ack(delivery_tag=method.delivery_tag)
            else:
                next_retry = retry_count + 1
                Logger.warning(
                    f"Message {data['id']} failed. Retrying ({next_retry}/{MAX_RETRIES})..."
                )

                # 🧩 Increment retry count in DB
                increment_retry_count(db, data["id"])  # type: ignore

                ch.basic_publish(
                    exchange="",
                    routing_key=RETRY_QUEUE,
                    body=body,
                    properties=pika.BasicProperties(
                        delivery_mode=2,
                        headers={"x-retry-count": next_retry},
                    ),
                )
                ch.basic_ack(delivery_tag=method.delivery_tag)

    except Exception as e:
        Logger.error(f"Error processing message: {e}")
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
    finally:
        db.close()


def start_consumer():
    while True:
        try:
            connection = pika.BlockingConnection(
                pika.ConnectionParameters(RABBITMQ_HOST)  # type: ignore
            )
            channel = connection.channel()
            setup_queues(channel)
            channel.basic_qos(prefetch_count=1)
            channel.basic_consume(queue=MAIN_QUEUE, on_message_callback=process_message)
            Logger.info("Worker started. Waiting for messages...")
            channel.start_consuming()
        except AMQPConnectionError as e:
            Logger.error(f"RabbitMQ connection error: {e}, retrying in 5s...")
            time.sleep(5)
        except ChannelClosedByBroker as e:
            Logger.error(f"Channel closed by broker: {e}, restarting consumer...")
            time.sleep(5)
        except KeyboardInterrupt:
            Logger.info("Worker stopped by user.")
            break


if __name__ == "__main__":
    start_consumer()
