import os
os.environ["PASSLIB_PURE"] = "1"
from datetime import datetime, timedelta
from jose import jwt, JWTError
from passlib.hash import bcrypt
from app.core.config import settings
from uuid import UUID

# Fix bcrypt backend issues in some Docker envs

MAX_BCRYPT_BYTES = 72

def hash_password(password: str) -> str:
    truncated_bytes = password.encode("utf-8")[:MAX_BCRYPT_BYTES]
    truncated_str = truncated_bytes.decode("utf-8", errors="ignore")
    return bcrypt.hash(truncated_str)

def verify_password(password: str, hashed_password: str) -> bool:
    truncated_bytes = password.encode("utf-8")[:MAX_BCRYPT_BYTES]
    truncated_str = truncated_bytes.decode("utf-8", errors="ignore")
    return bcrypt.verify(truncated_str, hashed_password)

def create_access_token(*, subject: str):
    to_encode = {"sub": subject}
    expire = datetime.utcnow() + timedelta(minutes=settings.access_token_expire_minutes)
    to_encode.update({"exp": expire})
    return jwt.encode(to_encode, settings.jwt_secret, algorithm=settings.jwt_algorithm)


def verify_token(token: str) -> UUID:
    """Verify JWT token and return user_id"""
    try:
        if token.startswith("Bearer "):
            token = token.split(" ")[1]

        payload = jwt.decode(token, settings.jwt_secret, algorithms=[settings.jwt_algorithm])
        user_id: str = payload.get("sub")

        if not user_id:
            raise JWTError("Invalid token payload")

        return UUID(user_id)

    except JWTError as e:
        raise Exception(f"Token verification failed: {str(e)}")
