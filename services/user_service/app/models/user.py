from sqlalchemy import Column, String, DateTime, func, JSON
from sqlalchemy.dialects.postgresql import UUID
import uuid
from app.core.db import Base


class User(Base):
    __tablename__ = "users"

    user_id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    email = Column(String(255), unique=True, nullable=False)
    full_name = Column(String(255))
    password_hash = Column(String(255), nullable=False)
    phone = Column(String(20))
    push_token = Column(String(255), nullable=True)
    preferences = Column(JSON, default={"email": True, "push": True})
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), onupdate=func.now())
