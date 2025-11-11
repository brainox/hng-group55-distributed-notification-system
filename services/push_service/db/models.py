#!/usr/bin/env python3
import enum

from sqlmodel import Column, Integer, String, Enum, DateTime, func

from .database import Base


class StatusEnum(str, enum.Enum):
    pending = "pending"
    success = "success"
    failed = "failed"


class PushMessage(Base):
    __tablename__ = "push_messages"

    id = Column(Integer, primary_key=True, index=True)
    title = Column(String)
    body = Column(String)
    token = Column(String)
    status = Column(Enum(StatusEnum), default=StatusEnum.pending)
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), onupdate=func.now())
