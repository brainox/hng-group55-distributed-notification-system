#!/usr/bin/env python3
import enum
from datetime import datetime, timezone

from pydantic import field_serializer
from sqlmodel import SQLModel, String, func, Field


class StatusEnum(str, enum.Enum):
    pending = "pending"
    success = "success"
    failed = "failed"


class PushMessage(SQLModel, table=True):
    __name__ = "push_messages"

    id: int = Field(default=None, primary_key=True)
    title: str = Field(String)
    body: str = Field(String)
    token: str = Field(String)
    status: StatusEnum = Field(default=StatusEnum.pending)
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(
        default_factory=lambda: datetime.now(timezone.utc),
        sa_column_kwargs={"onupdate": func.now()},
    )
    retry_count: int = Field(default=0)

    @field_serializer("created_at", "updated_at", when_used="always")
    def serialize_datetime(self, dt: datetime) -> str:
        return dt.isoformat()
