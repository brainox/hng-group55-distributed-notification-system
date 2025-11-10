from datetime import datetime, timedelta
from jose import jwt
from passlib.context import CryptContext
from app.core.config import settings
import os
os.environ["PASSLIB_PURE"] = "1"
from passlib.hash import bcrypt


MAX_BCRYPT_BYTES = 72

def hash_password(password: str) -> str:
    """
    Hash the password using bcrypt, safely truncated to 72 bytes.
    """
    # Convert to bytes and truncate safely
    truncated_bytes = password.encode("utf-8")[:MAX_BCRYPT_BYTES]
    truncated_str = truncated_bytes.decode("utf-8", errors="ignore")
    return bcrypt.hash(truncated_str)

def verify_password(password: str, hashed_password: str) -> bool:
    """
    Verify password after truncating to 72 bytes.
    """
    truncated_bytes = password.encode("utf-8")[:MAX_BCRYPT_BYTES]
    truncated_str = truncated_bytes.decode("utf-8", errors="ignore")
    return bcrypt.verify(truncated_str, hashed_password)


def create_access_token(data: dict):
    to_encode = data.copy()
    expire = datetime.utcnow() + timedelta(minutes=settings.access_token_expire_minutes)
    to_encode.update({"exp": expire})
    return jwt.encode(to_encode, settings.jwt_secret, algorithm=settings.jwt_algorithm)
