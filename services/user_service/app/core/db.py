from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker, declarative_base
import redis.asyncio as redis

from app.core.config import settings

Base = declarative_base()

# --- SQLAlchemy setup ---
engine = create_async_engine(settings.database_url, echo=False, future=True)
async_session = sessionmaker(engine, expire_on_commit=False, class_=AsyncSession)

# --- Redis setup ---
redis_client = None  # use a different name to avoid shadowing the module


async def init_db():
    """Initialize the database connection."""
    async with engine.begin() as conn:
        pass  # You can create tables here if needed


async def init_redis():
    """Initialize Redis connection (async)."""
    global redis_client
    redis_client = redis.from_url(
        settings.redis_url,
        decode_responses=True
    )
    # Optional: verify connection
    await redis_client.ping()

