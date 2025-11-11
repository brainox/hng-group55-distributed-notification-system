"""Simple Redis-backed queue helper. Replace with your preferred queue (RabbitMQ, SQS, etc.)"""

import json
from typing import Dict, Optional, Any

import pika
import redis


class RabbitMQQueue:
    def __init__(
        self,
        amqp_url: str = "amqp://guest:guest@rabbitmq:45672/%2f",
        queue_name: str = "push_queue",
    ):
        if pika is None:
            raise RuntimeError("pika not installed; install pika to use RabbitMQQueue")
        self.queue_name = queue_name
        self.conn = pika.BlockingConnection(
            pika.ConnectionParameters("localhost")
        )
        self.channel = self.conn.channel()
        self.channel.queue_declare(queue=self.queue_name, durable=True)
        self.channel.queue_declare(queue=self.queue_name + "_dlq", durable=True)

    def push(self, message: Dict[str, Any]) -> None:
        body = json.dumps(message)
        self.channel.basic_publish(
            exchange="",
            routing_key=self.queue_name,
            body=body,
            properties=pika.BasicProperties(delivery_mode=2),
        )

    def pop(self, block: bool = True, timeout: int = 5) -> Optional[Dict[str, Any]]:
        """Use basic_get polling. If found, return payload plus `_rmq_delivery_tag` to allow ack."""
        import time

        waited = 0.0
        while True:
            method_frame, header_frame, body = self.channel.basic_get(
                self.queue_name, auto_ack=False
            )
            if method_frame:
                try:
                    payload = json.loads(body)
                except Exception:
                    self.channel.basic_ack(method_frame.delivery_tag)
                    return None
                payload["_rmq_delivery_tag"] = method_frame.delivery_tag
                return payload
            if not block:
                return None
            time.sleep(0.5)
            waited += 0.5
            if timeout and waited >= timeout:
                return None

    def ack(self, delivery_tag):
        self.channel.basic_ack(delivery_tag)

    def nack(self, delivery_tag, requeue=False):
        self.channel.basic_nack(delivery_tag, requeue=requeue)

    def push_dlq(self, message: Dict[str, Any]) -> None:
        body = json.dumps(message)
        self.channel.basic_publish(
            exchange="",
            routing_key=self.queue_name + "_dlq",
            body=body,
            properties=pika.BasicProperties(delivery_mode=2),
        )
    
    def cleanup(self) -> None:
        self.channel.close()
