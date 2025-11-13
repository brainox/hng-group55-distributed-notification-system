#!/usr/bin/env python3

from sqlmodel import Session

from db import models
from logger import Logger


def increment_retry_count(db: Session, message_id: int):
    """Increment retry count for a given message ID."""
    msg = (
        db.query(models.PushMessage).filter(models.PushMessage.id == message_id).first()  # type: ignore
    )
    if msg:
        msg.retry_count = msg.retry_count + 1
        db.commit()
        Logger.info(f"Retry count for message {msg.id} -> {msg.retry_count}")
    else:
        Logger.warning(f"Message {message_id} not found for retry increment")
