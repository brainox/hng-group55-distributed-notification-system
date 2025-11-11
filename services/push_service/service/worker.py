import re
from typing import Dict, Any

from .senders import FCMSender
from .queue import RabbitMQQueue

from logger import Logger


class PushWorker:
    def __init__(
        self,
        queue: RabbitMQQueue,
        fcm_sender: FCMSender = None,  # type: ignore
    ):
        self.queue = queue
        self.fcm_sender = fcm_sender

    def handle_message(self, msg: Dict[str, Any]) -> None:
        """Message format (example):
        {
          "provider": "fcm" | "onesignal" | "webpush",
          "target": "token_or_player_id" or subscription dict for webpush,
          "title": "...",
          "body": "...",
          "image": "https://...",
          "url": "https://...",
          "data": {...}
        }
        """
        title = msg.get("title", "")
        body = msg.get("body", "")
        image = msg.get("image")
        url = msg.get("url")
        data = msg.get("data")

        token = msg.get("target")
        if not self.validate_fcm_token(str(token)):
            Logger.warning("Invalid FCM token: %s", token)
            return
        if not self.fcm_sender:
            Logger.error("FCM sender not configured")
            return
        resp = self.fcm_sender.send(
            str(token), title, body, image=image, data=data, click_action=url
        )
        Logger.info("FCM response: %s", resp)

    def validate_fcm_token(self, token: str) -> bool:
        return bool(token) and bool(re.fullmatch(r"[A-Za-z0-9:_\-]{20,400}", token))

    async def run_forever(self, poll_timeout: int = 5):
        Logger.info("Starting push worker loop...")
        while True:
            msg = self.queue.pop(block=True, timeout=poll_timeout)
            if msg is None:
                continue
            try:
                self.handle_message(msg)
            except Exception as e:
                Logger.exception("Failed to handle message: %s", e)
