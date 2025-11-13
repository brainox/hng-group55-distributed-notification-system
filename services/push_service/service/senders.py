#!/usr/bin/env python3

import time
import json
import base64
import functools
from typing import Optional, Dict, Any

from firebase_admin import messaging, initialize_app, credentials

from logger import Logger


def retry(max_retries=3, delay=2):
    """
    Retry a function up to `max_retries` times with a `delay` (in seconds)
    between attempts.

    Usage:
        @retry(max_retries=3, delay=1)
        def unreliable_func():
            ...
    """

    def decorator(func):
        @functools.wraps(func)
        def wrapper(*args, **kwargs):
            for attempt in range(1, max_retries + 1):
                try:
                    return func(*args, **kwargs)
                except Exception as e:
                    Logger.warning(
                        f"[Retry {attempt}/{max_retries}] {func.__name__} failed: {e}"
                    )
                    if attempt < max_retries:
                        time.sleep(delay)
                    else:
                        Logger.error(f"All {max_retries} retries failed.")
                        raise

        return wrapper

    return decorator


class FCMSender:
    def __init__(self, credential_path: str):
        self.credential_path = credential_path
        self.cred = credentials.Certificate(credential_path)
        initialize_app(self.cred)
        Logger.info("Initializing FCMSEnder complete")

    @retry(max_retries=3, delay=2)
    def send(
        self,
        token: str,
        title: str,
        body: str,
        image: Optional[str] = None,
        data: Optional[Dict[str, Any]] = None,
        click_action: Optional[str] = None,
        **kwargs,
    ) -> bool:
        # Using firebase_admin SDK for sending
        try:
            payload = {
                "title": title,
                "body": body,
            }
            if image:
                payload["image"] = image
            if click_action:
                payload["click_action"] = click_action
            notification = messaging.Notification(**payload)
            message = messaging.Message(
                notification=notification,
                data=data,
                token=token,
            )

            messaging.send(message)
            Logger.info("FCM message sent to token: %s", token)
            return True
        except Exception as e:
            return False
