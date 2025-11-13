from sqlalchemy import Column, String, ForeignKey, DateTime, func
from sqlalchemy.dialects.postgresql import UUID
import uuid
from app.core.db import Base

class PushToken(Base):
    __tablename__ = "push_tokens"

    token_id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    user_id = Column(UUID(as_uuid=True), ForeignKey("users.user_id", ondelete="CASCADE"))
    token = Column(String, nullable=False)
    provider = Column(String(50))  # fcm,
    created_at = Column(DateTime(timezone=True), server_default=func.now())
