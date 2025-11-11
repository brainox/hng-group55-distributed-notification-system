import pika
import json
import logging
from sqlalchemy.orm import Session
from .database import SessionLocal
from . import models
import random
import time

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("consumer")

def process_message(ch, method, properties, body):
    db: Session = SessionLocal()
    try:
        data = json.loads(body)
        logger.info(f"Processing message {data['id']}: {data['title']}")

        # simulate push sending (replace with FCM/OneSignal call)
        time.sleep(2)
        success = random.choice([True, True, False]) 

        msg = db.query(models.PushMessage).filter(models.PushMessage.id == data["id"]).first()
        if msg:
            msg.status = models.StatusEnum.success if success else models.StatusEnum.failed
            db.commit()
            logger.info(f"Message {msg.id} updated -> {msg.status}")
        else:
            logger.warning(f"Message {data['id']} not found in DB")

        ch.basic_ack(delivery_tag=method.delivery_tag)
    except Exception as e:
        logger.error(f"Error processing message: {e}")
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
    finally:
        db.close()

def start_consumer():
    while True:
        try:
            connection = pika.BlockingConnection(pika.ConnectionParameters(host="rabbitmq"))
            channel = connection.channel()
            channel.queue_declare(queue="push_queue", durable=True)
            channel.basic_qos(prefetch_count=1)
            channel.basic_consume(queue="push_queue", on_message_callback=process_message)
            logger.info("Worker started. Waiting for messages...")
            channel.start_consuming()
        except pika.exceptions.AMQPConnectionError as e:
            logger.error(f"RabbitMQ connection error: {e}, retrying in 5s")
            time.sleep(5)
