#!/usr/bin/env python3

import json
import base64
from typing import Optional, Dict, Any

from firebase_admin import messaging, initialize_app, credentials

from logger import Logger


class FCMSender:
    def __init__(self, credential_path: str):
        self.credential_path = credential_path
        self.cred = credentials.Certificate(credential_path)
        initialize_app(self.cred)
        Logger.info("Initializing FCMSEnder complete")

    async def send(
        self,
        token: str,
        title: str,
        body: str,
        image: Optional[str] = None,
        data: Optional[Dict[str, Any]] = None,
        click_action: Optional[str] = None,
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

            await messaging.send_each_async([message])
            Logger.info("FCM message sent to token: %s", token)
            return True
        except Exception as e:
            return False
